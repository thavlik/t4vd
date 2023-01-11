import 'package:t4vd/model.dart';
import 'package:t4vd/sources/video_details.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';

import '../api.dart';
import '../context.dart';

class InputVideosPage extends StatefulWidget {
  const InputVideosPage({super.key});

  @override
  State<InputVideosPage> createState() => InputVideosPageState();
}

class PendingVideoListItem extends StatelessWidget {
  const PendingVideoListItem({
    super.key,
    required this.id,
    required this.message,
    this.title,
    this.showProgressIndicator = false,
    this.thumbnail = false,
  });

  final String id;
  final String? title;
  final String message;
  final bool showProgressIndicator;
  final bool thumbnail;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: () {},
      child: Container(
        decoration: BoxDecoration(
          border: Border(
            bottom: BorderSide(
              color: Theme.of(context).dividerColor,
            ),
          ),
        ),
        height: 100,
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              AspectRatio(
                aspectRatio: 1280.0 / 720.0,
                child: Stack(
                  children: [
                    Container(
                      decoration: BoxDecoration(
                        color: Colors.white,
                        borderRadius:
                            const BorderRadius.all(Radius.circular(8.0)),
                        image: thumbnail
                            ? DecorationImage(
                                image: NetworkImage(videoThumbnail(id)),
                                alignment: const Alignment(0, 0),
                                fit: BoxFit.cover,
                              )
                            : null,
                      ),
                    ),
                    Center(
                      child: showProgressIndicator
                          ? const CircularProgressIndicator()
                          : const Icon(
                              Icons.question_mark,
                              color: Colors.black,
                            ),
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 8),
              Expanded(
                child: Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        title ?? id,
                        style: Theme.of(context).textTheme.titleLarge,
                        overflow: TextOverflow.ellipsis,
                      ),
                      const SizedBox(height: 4),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.start,
                        children: [
                          Text(
                            message,
                            style: Theme.of(context).textTheme.bodySmall,
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class InputVideoListItem extends StatefulWidget {
  const InputVideoListItem({
    super.key,
    required this.id,
    this.info,
    this.editMode = false,
    this.onDelete,
  });

  final String id;
  final VideoInfo? info;
  final bool editMode;
  final void Function()? onDelete;

  @override
  State<InputVideoListItem> createState() => _InputVideoListItemState();
}

class _InputVideoListItemState extends State<InputVideoListItem> {
  Offset? lastPosition;

  @override
  Widget build(BuildContext context) {
    final info = widget.info;
    final secs = 777; //model.duration;
    final m = (secs.toDouble() / 60.0).floor();
    final s = secs - m * 60;
    final duration = "${m}m${s}s";
    return InkWell(
      key: Key('video-${widget.id}'),
      onTap: info != null
          ? () => Navigator.of(context).push(MaterialPageRoute(
              builder: (context) => VideoDetailsPage(Video(
                    id: widget.id,
                    info: info,
                  ))))
          : () => {},
      onLongPress: () {
        if (lastPosition != null) {
          showVideoContextMenu(
            context,
            lastPosition!,
            id: widget.id,
            onBlacklist: () {},
          );
        }
      },
      child: Listener(
        onPointerDown: (PointerDownEvent event) {
          lastPosition = event.position;
          if (event.kind != PointerDeviceKind.mouse ||
              event.buttons != kSecondaryMouseButton) {
            return;
          }
          showVideoContextMenu(
            context,
            event.position,
            id: widget.id,
            onBlacklist: () {},
          );
        },
        child: Container(
          height: 100,
          decoration: BoxDecoration(
            border: Border(
              bottom: BorderSide(
                color: Theme.of(context).dividerColor,
              ),
            ),
          ),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.start,
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                Align(
                  child: AspectRatio(
                    aspectRatio: 1280.0 / 720.0,
                    child: Container(
                      decoration: BoxDecoration(
                        borderRadius:
                            const BorderRadius.all(Radius.circular(8.0)),
                        image: DecorationImage(
                          image: NetworkImage(videoThumbnail(widget.id)),
                          alignment: const Alignment(0, 0),
                          fit: BoxFit.cover,
                        ),
                      ),
                    ),
                  ),
                ),
                const SizedBox(width: 8),
                Expanded(
                  flex: 1,
                  child: Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          info?.title ?? widget.id,
                          style: Theme.of(context).textTheme.titleLarge,
                          overflow: TextOverflow.ellipsis,
                        ),
                        Text(
                          info != null
                              ? "${info.channel} • $duration"
                              : "<unknown channel>",
                          style: Theme.of(context).textTheme.bodySmall,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ],
                    ),
                  ),
                ),
                Visibility(
                  visible: widget.editMode,
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(16, 0, 0, 0),
                    child: IconButton(
                      onPressed: widget.onDelete,
                      icon: Icon(
                        Icons.remove_circle,
                        color: Theme.of(context).buttonTheme.colorScheme!.error,
                      ),
                    ),
                  ),
                ),
                Visibility(
                  visible: !widget.editMode,
                  child: GestureDetector(
                    onTapDown: (details) {
                      showVideoContextMenu(
                        context,
                        details.globalPosition,
                        id: widget.id,
                        onBlacklist: () {},
                      );
                    },
                    child: Padding(
                      padding: const EdgeInsets.fromLTRB(16, 0, 0, 0),
                      child: Icon(
                        Icons.more_vert,
                        color:
                            Theme.of(context).buttonTheme.colorScheme!.primary,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

class AddVideoResult {
  final String urlOrId;
  final bool blacklist;

  AddVideoResult({
    required this.urlOrId,
    required this.blacklist,
  });
}

class InputVideosPageState extends State<InputVideosPage> {
  bool editMode = false;
  bool loading = true;

  final addVideoController = TextEditingController();
  final focusNode = FocusNode();

  @override
  void initState() {
    super.initState();
    ScopedModel.of<BJJModel>(context)
        .refreshVideos(Navigator.of(context))
        .then((value) => setState(() {
              loading = false;
            }));
  }

  void addVideo(
    BuildContext context,
    String value,
    bool blacklist,
  ) async {
    await ScopedModel.of<BJJModel>(context).addVideo(
      nav: Navigator.of(context),
      input: value,
      blacklist: blacklist,
    );
  }

  void deleteVideo(BuildContext context, String id, bool blacklist) async {
    await ScopedModel.of<BJJModel>(context).removeVideo(
      nav: Navigator.of(context),
      id: id,
      blacklist: blacklist,
    );
  }

  void showAddVideoDialog(BuildContext context) async {
    addVideoController.clear();
    final AddVideoResult? result = await showDialog(
        context: context,
        builder: (context) {
          bool blacklist = false;
          return StatefulBuilder(
            builder: (context, setState) => AlertDialog(
              actions: [
                TextButton(
                  onPressed: () => Navigator.of(context).pop(),
                  child: Text(
                    "Cancel",
                    style: Theme.of(context).textTheme.labelLarge,
                  ),
                ),
                TextButton(
                  key: const Key('confirmAddVideo'),
                  onPressed: () {
                    if (addVideoController.text.isEmpty) {
                      focusNode.requestFocus();
                      return;
                    }
                    Navigator.of(context).pop(AddVideoResult(
                      urlOrId: addVideoController.text,
                      blacklist: blacklist,
                    ));
                  },
                  child: Text(
                    "Confirm",
                    style: Theme.of(context).textTheme.labelLarge,
                  ),
                ),
              ],
              title: Text(
                "Add Video",
                style: Theme.of(context).textTheme.titleLarge,
              ),
              content: SizedBox(
                height: 110,
                child: Column(
                  children: [
                    TextField(
                      key: const Key('addVideoInput'),
                      decoration: const InputDecoration(
                        labelText: "Please enter a video URL or ID.",
                      ),
                      focusNode: focusNode,
                      autofocus: true,
                      controller: addVideoController,
                    ),
                    const SizedBox(height: 16),
                    Row(
                      children: [
                        Checkbox(
                            value: blacklist,
                            onChanged: (bool? value) {
                              setState(() {
                                blacklist = value ?? false;
                              });
                            }),
                        Text(
                          "Blacklist",
                          style: Theme.of(context).textTheme.bodySmall,
                        )
                      ],
                    ),
                  ],
                ),
              ),
            ),
          );
        });
    if (!mounted || result == null) return;
    addVideo(context, result.urlOrId, result.blacklist);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          "Input Videos",
          style: Theme.of(context).textTheme.titleLarge,
        ),
        leading: IconButton(
          key: const Key('videosNavBack'),
          icon: const Icon(Icons.navigate_before),
          onPressed: () => Navigator.of(context).pop(),
        ),
        actions: [
          IconButton(
            onPressed: () {
              setState(() {
                editMode = !editMode;
              });
            },
            icon: Icon(editMode ? Icons.cancel : Icons.edit),
          )
        ],
      ),
      body: Stack(
        children: [
          ScopedModelDescendant<BJJModel>(builder: (context, child, model) {
            return ListView(
              children: [
                const PendingVideoListItem(
                  id: 'ExqT2SwW1qQ',
                  message: 'Tyler Spangler • 320 MiB downloaded • 1.2 MiB/sec',
                  showProgressIndicator: true,
                  title:
                      'I Survived The Highest Rated Jiu Jitsu Gyms In Las Vegas',
                  thumbnail: true,
                ),
                const PendingVideoListItem(
                  id: 'hWg3ooN2ia0',
                  message: 'Tyler Spangler • Pending download • #2 in queue',
                  title: 'The Last Video Of Bath Salts Ben',
                  thumbnail: true,
                ),
                const PendingVideoListItem(
                  id: 'YZXXhpjCXiE',
                  message: 'Querying • 4s elapsed',
                ),
                const PendingVideoListItem(
                  id: 'morD58OZmy0',
                  message: 'Pending query • #2 in queue',
                ),
                ...model.videos
                    .where((vid) => !vid.blacklist)
                    .map((vid) => InputVideoListItem(
                          editMode: editMode,
                          onDelete: () => deleteVideo(context, vid.id, false),
                          id: vid.id,
                          info: vid.info,
                        ))
                    .toList(),
                Visibility(
                    visible:
                        model.videos.indexWhere((ch) => !ch.blacklist) == -1 &&
                            !loading,
                    child: Padding(
                      padding: const EdgeInsets.all(16.0),
                      child: Opacity(
                        opacity: 0.5,
                        child: Text(
                          "There are no input videos. Add one by pressing the + sign in the bottom right corner.",
                          style: Theme.of(context).textTheme.bodyMedium,
                        ),
                      ),
                    )),
                Container(
                  decoration: BoxDecoration(
                    border: Border(
                        bottom: BorderSide(
                      color: Theme.of(context).dividerColor,
                    )),
                  ),
                  child: Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Text(
                      "Blacklisted Videos",
                      style: Theme.of(context).textTheme.titleLarge,
                    ),
                  ),
                ),
                ...model.videos
                    .where((vid) => vid.blacklist)
                    .map((vid) => InputVideoListItem(
                          editMode: editMode,
                          onDelete: () => deleteVideo(context, vid.id, true),
                          id: vid.id,
                          info: vid.info,
                        ))
                    .toList(),
                Visibility(
                    visible:
                        model.videos.indexWhere((ch) => ch.blacklist) == -1 &&
                            !loading,
                    child: Padding(
                      padding: const EdgeInsets.all(16.0),
                      child: Opacity(
                        opacity: 0.5,
                        child: Text(
                          "There are no blacklisted videos.",
                          style: Theme.of(context).textTheme.bodyMedium,
                        ),
                      ),
                    )),
              ],
            );
          }),
          Positioned(
              bottom: 16,
              right: 16,
              child: FloatingActionButton(
                key: const Key('addVideo'),
                onPressed: () => showAddVideoDialog(context),
                child: const Icon(Icons.add),
              )),
          Visibility(
            visible: loading,
            child: const Center(
              child: CircularProgressIndicator(),
            ),
          ),
        ],
      ),
    );
  }
}
