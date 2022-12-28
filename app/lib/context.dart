import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';

Future<void> showVideoContextMenu(
  BuildContext context,
  Offset globalPosition, {
  required String id,
  required void Function() onBlacklist,
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
                  color: Theme.of(context).textTheme.button!.color,
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
              const Icon(Icons.delete),
              const SizedBox(width: 4),
              Text(
                "Blacklist Video",
                style: TextStyle(
                  color: Theme.of(context).textTheme.button!.color,
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
  } else if (result == 1) {
    final url = Uri.parse("https://youtube.com/watch?v=$id");
    if (!await launchUrl(url)) {
      throw 'Could not launch $url';
    }
  } else if (result == 2) {
    throw 'unimplemented';
    //await blacklistVideo(id);
    //onBlacklist();
  } else {
    throw ErrorSummary("unreachable branch detected");
  }
}

Future<void> showChannelContextMenu(
  BuildContext context,
  Offset globalPosition, {
  required String id,
  required void Function() onBlacklist,
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
        value: 1,
        child: Padding(
          padding: const EdgeInsets.only(left: 0, right: 0),
          child: Row(
            children: [
              const Icon(Icons.open_in_browser),
              const SizedBox(width: 4),
              Text(
                "View on YouTube",
                style: TextStyle(
                  color: Theme.of(context).textTheme.button!.color,
                ),
              ),
            ],
          ),
        ),
      ),
      PopupMenuItem(
        value: 2,
        child: Padding(
          padding: const EdgeInsets.only(left: 0, right: 0),
          child: Row(
            children: [
              const Icon(Icons.delete),
              const SizedBox(width: 4),
              Text(
                "Blacklist Video",
                style: TextStyle(
                  color: Theme.of(context).textTheme.button!.color,
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
  } else if (result == 1) {
    final url = Uri.parse("https://youtube.com/$id");
    if (!await launchUrl(url)) {
      throw 'Could not launch $url';
    }
  } else if (result == 2) {
    throw 'unimplemented';
    //await blacklistChannel(id);
    //onBlacklist();
  } else {
    throw ErrorSummary("unreachable branch detected");
  }
}

Future<void> showPlaylistContextMenu(
  BuildContext context,
  Offset globalPosition, {
  required String id,
  void Function()? onBlacklist,
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
        value: 1,
        child: Padding(
          padding: const EdgeInsets.only(left: 0, right: 40),
          child: Row(
            children: [
              const Icon(Icons.web),
              const SizedBox(width: 4),
              Text(
                "View on YouTube",
                style: TextStyle(
                  color: Theme.of(context).textTheme.button!.color,
                ),
              ),
            ],
          ),
        ),
      ),
      ...onBlacklist != null
          ? [
              PopupMenuItem(
                value: 2,
                child: Padding(
                  padding: const EdgeInsets.only(left: 0, right: 40),
                  child: Row(
                    children: [
                      const Icon(Icons.delete),
                      const SizedBox(width: 4),
                      Text(
                        "Blacklist Video",
                        style: TextStyle(
                          color: Theme.of(context).textTheme.button!.color,
                        ),
                      ),
                    ],
                  ),
                ),
              )
            ]
          : [],
    ],
    elevation: 8.0,
  );
  if (result == null) {
    return;
  } else if (result == 1) {
    final url = Uri.parse("https://youtube.com/watch?list=$id");
    if (!await launchUrl(url)) {
      throw 'Could not launch $url';
    }
  } else if (result == 2) {
    if (onBlacklist != null) {
      throw 'unimplemented';
      //await blacklistPlaylist(id);
      //onBlacklist();
    }
  } else {
    throw ErrorSummary("unreachable branch detected");
  }
}
