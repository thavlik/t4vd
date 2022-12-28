package server

import (
	"context"
	"time"

	"github.com/pkg/errors"
	slideshow "github.com/thavlik/t4vd/slideshow/pkg/api"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/scheduler"
	"github.com/thavlik/t4vd/compiler/pkg/api"
	"github.com/thavlik/t4vd/compiler/pkg/compiler"
	"github.com/thavlik/t4vd/compiler/pkg/datastore"
	seer "github.com/thavlik/t4vd/seer/pkg/api"
	sources "github.com/thavlik/t4vd/sources/pkg/api"
	"go.uber.org/zap"
)

func Entry(
	port int,
	sched scheduler.Scheduler,
	ds datastore.DataStore,
	sourcesClient sources.Sources,
	seerOpts base.ServiceOptions,
	slideshow slideshow.SlideShow,
	saveInterval time.Duration,
	compileOnStart bool,
	concurrency int,
	log *zap.Logger,
) error {
	s := NewServer(
		ds,
		sched,
		sourcesClient,
		seerOpts,
		slideshow,
		saveInterval,
		log,
	)
	pop := make(chan string)
	stopPopper := make(chan struct{}, 1)
	stoppedPopper := make(chan struct{})
	go func() {
		popper(
			sched,
			pop,
			stopPopper,
			log,
		)
		stoppedPopper <- struct{}{}
	}()
	stoppedCompile := make([]chan struct{}, concurrency)
	for i := 0; i < concurrency; i++ {
		stopped := make(chan struct{})
		stoppedCompile[i] = stopped
		go func(stopped chan<- struct{}) {
			compile(
				sched,
				ds,
				sourcesClient,
				seer.NewSeerClientFromOptions(seerOpts),
				pop,
				saveInterval,
				log,
			)
			stopped <- struct{}{}
		}(stopped)
	}
	defer func() {
		stopPopper <- struct{}{}
		<-stoppedPopper
		for _, stopped := range stoppedCompile {
			<-stopped
		}
	}()
	if compileOnStart {
		resp, err := sourcesClient.ListProjects(
			context.Background(),
			sources.ListProjectsRequest{})
		if err != nil {
			return errors.Wrap(err, "sources.ListProjects")
		}
		for _, project := range resp.Projects {
			go s.Compile(context.Background(), api.Compile{
				ProjectID: project.ID,
			})
		}
	}
	base.SignalReady(log)
	return s.ListenAndServe(port)
}

func popper(
	sched scheduler.Scheduler,
	pop chan<- string,
	stop <-chan struct{},
	log *zap.Logger,
) {
	notification := sched.Notify()
	defer close(pop)
	delay := 12 * time.Second
	for {
		start := time.Now()
		projectIDs, err := sched.List()
		if err != nil {
			panic(errors.Wrap(err, "scheduler.List"))
		}
		if len(projectIDs) > 0 {
			log.Debug("checking projects", zap.Strings("projectIDs", projectIDs))
			for _, projectID := range projectIDs {
				pop <- projectID
			}
		}
		remaining := delay - time.Since(start)
		if remaining > 0 {
			select {
			case <-stop:
				return
			case <-notification:
				continue
			case <-time.After(remaining):
				continue
			}
		}
	}
}

func compile(
	sched scheduler.Scheduler,
	ds datastore.DataStore,
	sourcesClient sources.Sources,
	seerClient seer.Seer,
	pop <-chan string,
	saveInterval time.Duration,
	log *zap.Logger,
) {
	// compile threads
	for {
		projectID, ok := <-pop
		if !ok {
			return
		}
		projectLog := log.With(zap.String("projectID", projectID))
		projectLog.Debug("locking project")
		start := time.Now()
		lock, err := sched.Lock(projectID)
		projectLog.Debug("locked project", base.Elapsed(start))
		if err == scheduler.ErrLocked {
			// go to the next project
			projectLog.Debug("project already locked")
			continue
		} else if err != nil {
			panic(errors.Wrap(err, "scheduler.Lock"))
		}
		projectLog.Debug("compiling project")
		func() {
			defer func() {
				if err := lock.Release(); err != nil {
					projectLog.Warn("error unlocking project",
						base.Elapsed(start),
						zap.Error(err))
				} else {
					projectLog.Debug("unlocked project", base.Elapsed(start))
				}
			}()
			stop := make(chan struct{}, 1)
			stopped := make(chan struct{})
			defer func() {
				stop <- struct{}{}
				<-stopped
			}()
			saved := make(chan *api.Dataset)
			onProgress := make(chan struct{}, 1)
			defer close(onProgress)
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				defer func() {
					cancel()
					stopped <- struct{}{}
				}()
				for {
					select {
					case <-stop:
						return
					case <-time.After(5 * time.Second):
						if err := lock.Extend(); err != nil {
							log.Warn("failed to extend lock",
								base.Elapsed(start),
								zap.Error(err))
							return
						}
						continue
					case _, ok := <-onProgress:
						if !ok {
							return
						}
						if err := lock.Extend(); err != nil {
							log.Warn("failed to extend lock",
								base.Elapsed(start),
								zap.Error(err))
							return
						}
					case _, ok := <-saved:
						if !ok {
							return
						}
						if err := lock.Extend(); err != nil {
							log.Warn("failed to extend lock",
								base.Elapsed(start),
								zap.Error(err))
							return
						}
					}
				}
			}()
			dataset, err := compiler.Compile(
				ctx,
				projectID,
				sourcesClient,
				seerClient,
				ds,
				saveInterval,
				saved,
				onProgress,
				log,
			)
			if err == compiler.ErrNoVideos {
				if err := sched.Remove(projectID); err != nil {
					projectLog.Warn("failed to remove project from scheduler", zap.Error(err))
				}
				return
			} else if err != nil {
				log.Error("error compiling dataset",
					base.Elapsed(start),
					zap.Error(err))
				return
			}
			log.Debug("compile process was successful")
			if err := lock.Extend(); err != nil {
				log.Warn("failed to extend lock",
					base.Elapsed(start),
					zap.Error(err))
				return
			}
			// make sure all the videos are cached
			videoIDs := make([]string, len(dataset.Videos))
			for i, video := range dataset.Videos {
				videoIDs[i] = video.ID
			}
			if _, err := seerClient.BulkScheduleVideoDownloads(
				context.Background(),
				seer.BulkScheduleVideoDownloads{
					VideoIDs: videoIDs,
				},
			); err != nil {
				projectLog.Warn("failed to add videos to scheduler", zap.Error(err))
				return
			}
			// remove project from scheduler
			if err := sched.Remove(projectID); err != nil {
				projectLog.Warn("failed to remove project from scheduler", zap.Error(err))
				return
			}
			projectLog.Debug("dataset compilation successful")
		}()
	}
}
