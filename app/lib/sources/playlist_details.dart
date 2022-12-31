import 'package:flutter/material.dart';

import '../api.dart';

class PlaylistDetailsPage extends StatelessWidget {
  const PlaylistDetailsPage(this.model, {super.key});

  final PlaylistListItem model;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          "Playlist Details",
          style: Theme.of(context).textTheme.headline6,
        ),
        actions: [
          IconButton(
            onPressed: () {},
            icon: const Icon(Icons.more_vert),
          ),
        ],
      ),
      body: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          infoSection(context, "ID", model.id),
          infoSection(context, "Title", model.title),
          infoSection(context, "Channel", model.channel),
          infoSection(context, "Channel ID", model.channelId),
          infoSection(context, "# Videos", model.numVideos.toString()),
        ],
      ),
    );
  }

  Padding infoSection(
    BuildContext context,
    String name,
    String value,
  ) {
    return Padding(
      padding: const EdgeInsets.all(8.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          Text(
            name,
            style: Theme.of(context).textTheme.headline6,
          ),
          const SizedBox(width: 16),
          Text(
            value,
            style: Theme.of(context).textTheme.bodyMedium,
          ),
        ],
      ),
    );
  }
}
