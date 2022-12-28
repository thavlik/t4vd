import 'package:bjjv/context.dart';
import 'package:bjjv/model.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';

import '../api.dart';

class InputPlaylistsPage extends StatefulWidget {
  const InputPlaylistsPage({super.key});

  @override
  State<InputPlaylistsPage> createState() => InputPlaylistsPageState();
}

class InputPlaylistListItem extends StatelessWidget {
  const InputPlaylistListItem({
    super.key,
    required this.editMode,
    required this.onDelete,
    required this.model,
  });

  final bool editMode;
  final void Function() onDelete;
  final PlaylistListItem model;

  @override
  Widget build(BuildContext context) {
    const secs = 0; //model.totalDuration;
    final m = (secs.toDouble() / 60.0).floor();
    final s = secs - m * 60;
    final duration = "${m}m${s}s";
    return InkWell(
      onTap: () {},
      child: Listener(
        onPointerDown: (PointerDownEvent event) {
          if (event.kind != PointerDeviceKind.mouse ||
              event.buttons != kSecondaryMouseButton) {
            return;
          }
          showPlaylistContextMenu(
            context,
            event.position,
            id: model.id,
            onBlacklist: () {
              // TODO: remove item from list
            },
          );
        },
        child: Container(
          decoration: BoxDecoration(
              border: Border(
                  bottom: BorderSide(color: Theme.of(context).dividerColor))),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                AspectRatio(
                  aspectRatio: 1280.0 / 720.0,
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius:
                          const BorderRadius.all(Radius.circular(8.0)),
                      image: DecorationImage(
                        image: NetworkImage(playlistThumbnail(model.id)),
                        alignment: const Alignment(0, 0),
                        fit: BoxFit.cover,
                      ),
                    ),
                  ),
                ),
                Expanded(
                  flex: 1,
                  child: Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          model.title,
                          style: Theme.of(context).textTheme.headline6,
                          overflow: TextOverflow.ellipsis,
                        ),
                        Text(
                          "${model.channel} • ${model.numVideos} videos • $duration total",
                          style: Theme.of(context).textTheme.bodySmall,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ],
                    ),
                  ),
                ),
                Visibility(
                  visible: editMode,
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(16, 0, 0, 0),
                    child: IconButton(
                      onPressed: onDelete,
                      icon: Icon(
                        Icons.remove_circle,
                        color: Theme.of(context).buttonTheme.colorScheme!.error,
                      ),
                    ),
                  ),
                ),
                Visibility(
                  visible: !editMode,
                  child: GestureDetector(
                    onTapDown: (details) {
                      showPlaylistContextMenu(
                        context,
                        details.globalPosition,
                        id: model.id,
                        onBlacklist: () {
                          // TODO: remove item from list
                        },
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

class AddPlaylistResult {
  final String urlOrId;
  final bool blacklist;

  AddPlaylistResult({
    required this.urlOrId,
    required this.blacklist,
  });
}

class InputPlaylistsPageState extends State<InputPlaylistsPage> {
  bool editMode = false;
  bool loading = true;

  final addPlaylistController = TextEditingController();
  final focusNode = FocusNode();

  @override
  void initState() {
    super.initState();
    ScopedModel.of<BJJModel>(context)
        .refreshPlaylists(context)
        .then((value) => setState(() => loading = false));
  }

  void addPlaylist(BuildContext context, String input, bool blacklist) async {
    await ScopedModel.of<BJJModel>(context).addPlaylist(
      context: context,
      input: input,
      blacklist: blacklist,
    );
  }

  void deletePlaylist(BuildContext context, String id, bool blacklist) async {
    await ScopedModel.of<BJJModel>(context).removePlaylist(
      context: context,
      id: id,
      blacklist: blacklist,
    );
  }

  void showAddPlaylistDialog(BuildContext context) async {
    addPlaylistController.clear();
    final AddPlaylistResult? result = await showDialog(
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
                    style: Theme.of(context).textTheme.button,
                  ),
                ),
                TextButton(
                  onPressed: () {
                    if (addPlaylistController.text.isEmpty) {
                      focusNode.requestFocus();
                      return;
                    }
                    Navigator.of(context).pop(AddPlaylistResult(
                      urlOrId: addPlaylistController.text,
                      blacklist: blacklist,
                    ));
                  },
                  child: Text(
                    "Confirm",
                    style: Theme.of(context).textTheme.button,
                  ),
                ),
              ],
              title: Text(
                "Add Playlist",
                style: Theme.of(context).textTheme.headline6,
              ),
              content: SizedBox(
                height: 110,
                child: Column(
                  children: [
                    TextField(
                      decoration: const InputDecoration(
                        labelText: "Please enter a playlist URL or ID.",
                      ),
                      focusNode: focusNode,
                      autofocus: true,
                      controller: addPlaylistController,
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
    if (result == null) {
      return;
    }
    if (!mounted) return;
    addPlaylist(context, result.urlOrId, result.blacklist);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          "Input Playlists",
          style: Theme.of(context).textTheme.headline6,
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
          ScopedModelDescendant<BJJModel>(
              builder: (context, child, model) => ListView(
                    children: [
                      ...model.playlists
                          .where((pl) => !pl.blacklist)
                          .map((pl) => InputPlaylistListItem(
                                editMode: editMode,
                                onDelete: () =>
                                    deletePlaylist(context, pl.id, false),
                                model: pl,
                              ))
                          .toList(),
                      Visibility(
                          visible: model.playlists
                                      .indexWhere((ch) => !ch.blacklist) ==
                                  -1 &&
                              !loading,
                          child: Padding(
                            padding: const EdgeInsets.all(16.0),
                            child: Opacity(
                              opacity: 0.5,
                              child: Text(
                                "There are no input playlists. Add one by pressing the + sign in the bottom right corner.",
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
                            "Blacklisted Playlists",
                            style: Theme.of(context).textTheme.headline6,
                          ),
                        ),
                      ),
                      ...model.playlists
                          .where((pl) => pl.blacklist)
                          .map((pl) => InputPlaylistListItem(
                                editMode: editMode,
                                onDelete: () =>
                                    deletePlaylist(context, pl.id, true),
                                model: pl,
                              ))
                          .toList(),
                      Visibility(
                          visible: model.playlists
                                      .indexWhere((ch) => ch.blacklist) ==
                                  -1 &&
                              !loading,
                          child: Padding(
                            padding: const EdgeInsets.all(16.0),
                            child: Opacity(
                              opacity: 0.5,
                              child: Text(
                                "There are no blacklisted playlists.",
                                style: Theme.of(context).textTheme.bodyMedium,
                              ),
                            ),
                          )),
                    ],
                  )),
          Positioned(
              bottom: 16,
              right: 16,
              child: FloatingActionButton(
                onPressed: () => showAddPlaylistDialog(context),
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
