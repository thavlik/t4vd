import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';
import 'package:url_launcher/url_launcher.dart';

import 'api.dart';
import 'model.dart';

class TagWidget extends StatefulWidget {
  const TagWidget(
    this.tag, {
    super.key,
    required this.onDelete,
  });

  final String tag;

  final void Function(String) onDelete;

  @override
  State<TagWidget> createState() => _TagWidgetState();
}

class _TagWidgetState extends State<TagWidget> {
  bool showDelete = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: (event) {
        setState(() {
          showDelete = true;
        });
      },
      onExit: (event) {
        setState(() {
          showDelete = false;
        });
      },
      child: Padding(
        padding: const EdgeInsets.all(4.0),
        child: Stack(
          clipBehavior: Clip.none,
          children: [
            Container(
              decoration: const BoxDecoration(
                color: Colors.pink,
                borderRadius: BorderRadius.all(Radius.circular(32)),
              ),
              child: Padding(
                padding: const EdgeInsets.all(8.0),
                child: Text(
                  widget.tag,
                  style: Theme.of(context).textTheme.bodySmall!.copyWith(
                        fontWeight: FontWeight.bold,
                        color: Colors.white,
                      ),
                ),
              ),
            ),
            ...showDelete
                ? [
                    Positioned(
                      top: -16,
                      right: -12,
                      child: IconButton(
                        onPressed: () => widget.onDelete(widget.tag),
                        padding: EdgeInsets.zero,
                        icon: Icon(
                          Icons.cancel_rounded,
                          color: Theme.of(context).primaryColor,
                        ),
                      ),
                    )
                  ]
                : [],
          ],
        ),
      ),
    );
  }
}

class TagsPage extends StatefulWidget {
  const TagsPage({
    super.key,
  });

  @override
  State<TagsPage> createState() => _TagsPageState();
}

class _TagsPageState extends State<TagsPage> {
  final _textController = TextEditingController();
  final _textNode = FocusNode();
  bool _loading = false;

  List<String> tags = [
    //"rear naked choke",
    //"side control",
    //"bottom armbar",
    //"guillotine",
    //"de la Riva",
    //"bicep slicer",
    //"kimura",
    //"heel hook",
    //"triangle",
    //"half guard",
    //"standing",
    //"russian tie-up",
    //"wrist lock",
    //"guard pull",
    //"full guard",
    //"anaconda",
    //"d'arce",
    //"kesa gatame",
    //"north south",
    //"lapel choke",
    //"ezequiel choke",
    //"gogoplata",
    //"paper cutter",
    //"americana",
    //"omoplata",
    //"spider guard",
    //"butterfly guard",
    //"50-50",
    //"knee-on-belly",
    //"lasso guard",
    //"boston crab",
    //"twister",
    //"kneebar",
    //"toe hold",
    //"straight ankle lock",
  ];

  @override
  void initState() {
    super.initState();
    final model = ScopedModel.of<BJJModel>(context);
    WidgetsBinding.instance.addPostFrameCallback((timeStamp) async {
      if (!mounted) return;
      if (model.markers.isEmpty || model.markerIndex == model.markers.length) {
        setState(() => _loading = true);
        try {
          await model.refreshMarkers(Navigator.of(context));
        } on InvalidCredentialsError {
          Navigator.of(context).pushNamed('/splash');
        } on ForbiddenError {
          Navigator.of(context).pushNamed('/splash');
        } finally {
          setState(() => _loading = false);
        }
      }
      if (!mounted) return;
      //model.precacheFrames(context);
    });
  }

  Future<void> skip(BuildContext context) async {
    if (_loading) return;
    await ScopedModel.of<BJJModel>(context).skip(Navigator.of(context));
  }

  Future<void> discard(BuildContext context) async {
    if (_loading) return;
    await ScopedModel.of<BJJModel>(context).discard(Navigator.of(context));
  }

  Future<void> submit(BuildContext context) async {
    if (_loading) return;
    await ScopedModel.of<BJJModel>(context).tag(
      nav: Navigator.of(context),
      tags: tags,
    );
  }

  void previous(BuildContext context) {
    if (_loading) return;
    final model = ScopedModel.of<BJJModel>(context);
    model.markerBack();
    setState(() => tags = model.currentMarker != null
        ? model.getTags(model.currentMarker!) ?? []
        : []);
  }

  void addTag(String tag) {
    if (_loading) return;
    setState(() {
      tags.add(tag);
      _textController.clear();
      _textNode.requestFocus();
    });
  }

  void removeTag(String tag) {
    if (_loading) return;
    setState(() => tags.remove(tag));
  }

  @override
  Widget build(BuildContext context) {
    final model = ScopedModel.of<BJJModel>(context);
    final currentMarker = model.currentMarker;
    return Stack(
      children: [
        Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            Expanded(
              child: Align(
                child: AspectRatio(
                  aspectRatio: 1920.0 / 1080.0,
                  child: Container(
                    decoration: BoxDecoration(
                      image: currentMarker != null
                          ? DecorationImage(
                              image: NetworkImage(currentMarker!.imageUrl),
                              alignment: Alignment(0, 0),
                              fit: BoxFit.cover,
                            )
                          : null,
                    ),
                    child: Padding(
                      padding: const EdgeInsets.all(8.0),
                      child: Wrap(
                        alignment: WrapAlignment.center,
                        crossAxisAlignment: WrapCrossAlignment.center,
                        runAlignment: WrapAlignment.end,
                        children: [
                          ...tags
                              .map((tag) => TagWidget(
                                    tag,
                                    onDelete: (tag) => removeTag(tag),
                                  ))
                              .toList()
                        ],
                      ),
                    ),
                  ),
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(0, 16, 0, 32),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  SizedBox(
                    width: 256,
                    child: TextField(
                      controller: _textController,
                      autofocus: true,
                      focusNode: _textNode,
                      textInputAction: TextInputAction.done,
                      onSubmitted: (value) => addTag(value),
                    ),
                  ),
                  IconButton(
                    onPressed: () => addTag(_textController.text),
                    icon: Icon(
                      Icons.send,
                      color: Theme.of(context).iconTheme.color,
                    ),
                  )
                ],
              ),
            ),
            const SizedBox(height: 48),
          ],
        ),
        Visibility(
          visible: ScopedModel.of<BJJModel>(context).markerIndex > 0,
          child: Positioned(
            top: 16,
            left: 16,
            child: FloatingActionButton(
              onPressed: () => previous(context),
              child: const Icon(Icons.navigate_before),
            ),
          ),
        ),
        Positioned(
          bottom: 16,
          left: 16,
          child: FloatingActionButton(
            onPressed: () => discard(context),
            child: const Icon(Icons.block),
          ),
        ),
        Positioned(
          bottom: 16,
          right: 16,
          child: FloatingActionButton(
            onPressed: () => submit(context),
            child: const Icon(Icons.done),
          ),
        ),
        if (currentMarker != null)
          Positioned(
            top: 0,
            right: 0,
            child: GestureDetector(
              onTapDown: (details) {
                showTagsContextMenu(
                  context,
                  details.globalPosition,
                  id: currentMarker.videoId,
                  startSeconds:
                      Duration(microseconds: currentMarker.time ~/ 1000)
                          .inSeconds,
                  onSkip: () => skip(context),
                  onDiscard: () => discard(context),
                );
              },
              child: Padding(
                padding: const EdgeInsets.all(8.0),
                child: Icon(
                  Icons.more_vert,
                  shadows: [
                    Shadow(
                      color: Colors.black.withAlpha(170),
                      blurRadius: 12.0,
                      offset: const Offset(0, 0),
                    )
                  ],
                  color: Colors.pink,
                ),
              ),
            ),
          ),
        if (_loading)
          const Center(
            child: CircularProgressIndicator(),
          ),
      ],
    );
  }
}

Future<void> showTagsContextMenu(
  BuildContext context,
  Offset globalPosition, {
  required String id,
  required int startSeconds,
  required void Function() onSkip,
  required void Function() onDiscard,
}) async {
  double left = globalPosition.dx;
  double top = globalPosition.dy;
  final result = await showMenu(
    color: Theme.of(context).cardColor,
    //add your color
    context: context,
    position: RelativeRect.fromLTRB(left, top, 0, 0),
    items: [
      PopupMenuItem(
        value: 0,
        child: Padding(
          padding: const EdgeInsets.only(left: 0, right: 40),
          child: Row(
            children: [
              const Icon(Icons.navigate_next),
              const SizedBox(width: 4),
              Text(
                "Skip",
                style: TextStyle(
                  color: Theme.of(context).textTheme.labelLarge!.color,
                ),
              ),
            ],
          ),
        ),
      ),
      PopupMenuItem(
        value: 1,
        child: Padding(
          padding: const EdgeInsets.only(left: 0, right: 40),
          child: Row(
            children: [
              const Icon(Icons.open_in_browser),
              const SizedBox(width: 4),
              Text(
                "View on YouTube",
                style: TextStyle(
                  color: Theme.of(context).textTheme.labelLarge!.color,
                ),
              ),
            ],
          ),
        ),
      ),
      PopupMenuItem(
        value: 2,
        child: Padding(
          padding: const EdgeInsets.only(left: 0, right: 40),
          child: Row(
            children: [
              const Icon(Icons.block),
              const SizedBox(width: 4),
              Text(
                "Discard Frame",
                style: TextStyle(
                  color: Theme.of(context).textTheme.labelLarge!.color,
                ),
              ),
            ],
          ),
        ),
      ),
    ],
    elevation: 8.0,
  );
  if (result == null) {
    return;
  } else if (result == 0) {
    onSkip();
  } else if (result == 1) {
    final url = Uri.parse("https://youtube.com/watch?v=$id&t=$startSeconds");
    if (!await launchUrl(url)) {
      throw 'Could not launch $url';
    }
  } else if (result == 2) {
    onDiscard();
  } else {
    throw ErrorSummary("unreachable branch detected");
  }
}
