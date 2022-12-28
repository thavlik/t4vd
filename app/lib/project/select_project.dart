import 'package:bjjv/model.dart';
import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';

import '../api.dart';

enum NavigatorBackBehavior {
  pop,
  signOut,
}

class SelectProjectPageArgs {
  final NavigatorBackBehavior? navigatorBackBehavior;

  SelectProjectPageArgs({
    this.navigatorBackBehavior,
  });
}

class SelectProjectPageArgsExtrator extends StatelessWidget {
  const SelectProjectPageArgsExtrator({super.key});

  @override
  Widget build(BuildContext context) {
    final args =
        ModalRoute.of(context)?.settings.arguments as SelectProjectPageArgs?;
    return SelectProjectPage(
      navigatorBackBehavior:
          args?.navigatorBackBehavior ?? NavigatorBackBehavior.signOut,
    );
  }
}

class SelectProjectPage extends StatefulWidget {
  const SelectProjectPage({
    super.key,
    required this.navigatorBackBehavior,
  });

  final NavigatorBackBehavior navigatorBackBehavior;

  @override
  State<SelectProjectPage> createState() => _SelectProjectPageState();
}

class _SelectProjectPageState extends State<SelectProjectPage> {
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((timeStamp) async {
      if (!mounted) return;
      final model = ScopedModel.of<BJJModel>(context);
      await model.refreshProjects(context);
      if (!mounted) return;
      setState(() => _loading = false);
    });
  }

  void onSelectProject(
    BuildContext context,
    Project p,
  ) async {
    await ScopedModel.of<BJJModel>(context).selectProject(
      context,
      p.id,
    );
    if (!mounted) return;
    Navigator.of(context).pop();
  }

  void onCreateProject(BuildContext context) async {
    await Navigator.of(context).pushNamed('/project/create');
    if (!mounted) return;
    if (ScopedModel.of<BJJModel>(context).hasProject) {
      Navigator.of(context).pop();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Select project'),
        automaticallyImplyLeading:
            widget.navigatorBackBehavior == NavigatorBackBehavior.pop,
        leading: widget.navigatorBackBehavior == NavigatorBackBehavior.signOut
            ? IconButton(
                key: const Key('selectProjectNavBack'),
                onPressed: () async {
                  if (widget.navigatorBackBehavior ==
                      NavigatorBackBehavior.pop) {
                    Navigator.of(context).pop();
                    return;
                  }
                  await ScopedModel.of<BJJModel>(context).signOut();
                  if (!mounted) return;
                  Navigator.of(context).pop();
                },
                icon: const Icon(Icons.navigate_before),
              )
            : null,
      ),
      body: ScopedModelDescendant<BJJModel>(
        builder: (context, child, model) => Center(
          child: Container(
            constraints: const BoxConstraints(
              maxWidth: 280,
            ),
            child: SingleChildScrollView(
              child: Padding(
                padding: const EdgeInsets.all(8.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    const SizedBox(height: 20),
                    Padding(
                      padding: const EdgeInsets.all(4.0),
                      child: Material(
                        color: Theme.of(context).backgroundColor,
                        child: InkWell(
                          key: const Key('createNewProject'),
                          onTap: () => onCreateProject(context),
                          child: SizedBox(
                            height: 64,
                            child: Padding(
                              padding: const EdgeInsets.all(16.0),
                              child: Center(
                                child: Text(
                                  "New project",
                                  style: Theme.of(context).textTheme.bodyLarge,
                                ),
                              ),
                            ),
                          ),
                        ),
                      ),
                    ),
                    ...model.projects
                        .map(
                          (p) => ProjectListItem(
                            p,
                            key: Key('project-${p.id}'),
                            onTap: () => onSelectProject(context, p),
                          ),
                        )
                        .toList(),
                    if (model.projects.isEmpty)
                      Padding(
                        padding: const EdgeInsets.all(4.0),
                        child: Opacity(
                          opacity: 0.7,
                          child: Container(
                            height: 64,
                            decoration: BoxDecoration(
                                border: Border.all(
                              width: 1,
                              color: Theme.of(context).dividerColor,
                            )),
                            child: Row(
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: [
                                _loading
                                    ? const Center(
                                        child: CircularProgressIndicator(),
                                      )
                                    : Flexible(
                                        child: Padding(
                                          padding: const EdgeInsets.all(20.0),
                                          child: Text(
                                            'There are no projects. Create one!',
                                            style: Theme.of(context)
                                                .textTheme
                                                .bodySmall,
                                          ),
                                        ),
                                      ),
                              ],
                            ),
                          ),
                        ),
                      ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class ProjectListItem extends StatelessWidget {
  const ProjectListItem(
    this.project, {
    required this.onTap,
    required super.key,
  });

  final Project project;
  final void Function() onTap;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(4.0),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          key: Key('selectProject-${project.name}'),
          onTap: onTap,
          child: Container(
            height: 64,
            decoration: BoxDecoration(
              color: Theme.of(context).backgroundColor.withAlpha(64),
              border: Border.all(
                width: 1.0,
                color: Theme.of(context).dividerColor,
              ),
            ),
            child: Padding(
              padding: const EdgeInsets.all(16.0),
              child: Center(
                child: Text(
                  project.name,
                  style: Theme.of(context).textTheme.bodyMedium,
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
