import 'package:t4vd/context.dart';
import 'package:t4vd/model.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';
import 'package:t4vd/sources/channel_details.dart';

import '../api.dart';

class InputChannelsPage extends StatefulWidget {
  const InputChannelsPage({super.key});

  @override
  State<InputChannelsPage> createState() => InputChannelsPageState();
}

class InputChannelListItem extends StatelessWidget {
  const InputChannelListItem({
    super.key,
    required this.editMode,
    required this.onDelete,
    required this.model,
  });

  final Channel model;
  final bool editMode;
  final void Function() onDelete;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      key: Key('channel-${model.id}'),
      onTap: () => Navigator.of(context).push(
          MaterialPageRoute(builder: (context) => ChannelDetailsPage(model))),
      child: Listener(
        onPointerDown: (PointerDownEvent event) {
          if (event.kind != PointerDeviceKind.mouse ||
              event.buttons != kSecondaryMouseButton) {
            return;
          }
          showChannelContextMenu(
            context,
            event.position,
            id: model.id,
            onBlacklist: () => onDelete(),
          );
        },
        child: Container(
          decoration: BoxDecoration(
              border: Border(
                  bottom: BorderSide(color: Theme.of(context).dividerColor))),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.start,
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                Container(
                  width: 48,
                  height: 48,
                  decoration: BoxDecoration(
                    borderRadius: const BorderRadius.all(Radius.circular(24.0)),
                    image: DecorationImage(
                      image: NetworkImage(channelAvatar(model.id)),
                      alignment: const Alignment(0, 0),
                      fit: BoxFit.cover,
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
                          model.info?.name ?? model.id,
                          style: Theme.of(context).textTheme.titleLarge,
                          overflow: TextOverflow.ellipsis,
                        ),
                        Text(
                          "#TODO_COUNT videos • #TODO_DURATION total",
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
                      showChannelContextMenu(
                        context,
                        details.globalPosition,
                        id: model.id,
                        onBlacklist: () => onDelete(),
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

class AddChannelResult {
  final String urlOrId;
  final bool blacklist;

  AddChannelResult({
    required this.urlOrId,
    required this.blacklist,
  });
}

class InputChannelsPageState extends State<InputChannelsPage> {
  bool editMode = false;
  bool loading = true;

  final focusNode = FocusNode();
  final addChannelController = TextEditingController();
  bool blacklistEnabled = false;

  @override
  void initState() {
    super.initState();
    ScopedModel.of<BJJModel>(context)
        .refreshChannels(Navigator.of(context))
        .then((value) {
      if (!mounted) return;
      setState(() => loading = false);
    });
  }

  void addChannel(BuildContext context, String input, bool blacklist) async {
    await ScopedModel.of<BJJModel>(context).addChannel(
      nav: Navigator.of(context),
      input: input,
      blacklist: blacklist,
    );
  }

  void deleteChannel(BuildContext context, String id, bool blacklist) async {
    await ScopedModel.of<BJJModel>(context).removeChannel(
      nav: Navigator.of(context),
      id: id,
      blacklist: blacklist,
    );
  }

  void showAddChannelDialog(BuildContext context) async {
    addChannelController.clear();
    final AddChannelResult? result = await showDialog(
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
                  key: const Key('confirmAddChannel'),
                  onPressed: () {
                    if (addChannelController.text.isEmpty) {
                      focusNode.requestFocus();
                      return;
                    }
                    Navigator.of(context).pop(AddChannelResult(
                      urlOrId: addChannelController.text,
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
                "Add Channel",
                style: Theme.of(context).textTheme.titleLarge,
              ),
              content: SizedBox(
                height: 110,
                child: Column(
                  children: [
                    TextField(
                      key: const Key('addChannelInput'),
                      decoration: const InputDecoration(
                        labelText: "Please enter a channel URL or ID.",
                      ),
                      focusNode: focusNode,
                      autofocus: true,
                      controller: addChannelController,
                    ),
                    const SizedBox(height: 16),
                    Row(
                      children: [
                        Checkbox(
                            value: blacklist,
                            onChanged: (bool? value) {
                              setState(() => blacklist = value ?? false);
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
    addChannel(context, result.urlOrId, result.blacklist);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          "Input Channels",
          style: Theme.of(context).textTheme.titleLarge,
        ),
        leading: IconButton(
          key: const Key('channelsNavBack'),
          icon: const Icon(Icons.navigate_before),
          onPressed: () => Navigator.of(context).pop(),
        ),
        actions: [
          IconButton(
            onPressed: () => setState(() => editMode = !editMode),
            icon: Icon(editMode ? Icons.cancel : Icons.edit),
          )
        ],
      ),
      body: Stack(
        children: [
          ScopedModelDescendant<BJJModel>(
            builder: (context, child, model) => ListView(
              children: [
                ...model.channels
                    .where((ch) => !ch.blacklist)
                    .map((ch) => InputChannelListItem(
                          editMode: editMode,
                          onDelete: () => deleteChannel(context, ch.id, false),
                          model: ch,
                        ))
                    .toList(),
                Visibility(
                    visible: model.channels.indexWhere((ch) => !ch.blacklist) ==
                            -1 &&
                        !loading,
                    child: Padding(
                      padding: const EdgeInsets.all(16.0),
                      child: Opacity(
                        opacity: 0.5,
                        child: Text(
                          "There are no input channels. Add one by pressing the + sign in the bottom right corner.",
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
                      "Blacklisted Channels",
                      style: Theme.of(context).textTheme.titleLarge,
                    ),
                  ),
                ),
                ...model.channels
                    .where((ch) => ch.blacklist)
                    .map((ch) => InputChannelListItem(
                          editMode: editMode,
                          onDelete: () => deleteChannel(context, ch.id, true),
                          model: ch,
                        ))
                    .toList(),
                Visibility(
                    visible:
                        model.channels.indexWhere((ch) => ch.blacklist) == -1 &&
                            !loading,
                    child: Padding(
                      padding: const EdgeInsets.all(16.0),
                      child: Opacity(
                        opacity: 0.5,
                        child: Text(
                          "There are no blacklisted channels.",
                          style: Theme.of(context).textTheme.bodyMedium,
                        ),
                      ),
                    )),
              ],
            ),
          ),
          Positioned(
              bottom: 16,
              right: 16,
              child: FloatingActionButton(
                key: const Key('addChannel'),
                onPressed: () => showAddChannelDialog(context),
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
