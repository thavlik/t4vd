import 'package:flutter/material.dart';

import '../api.dart';

class VideoDetailsPage extends StatelessWidget {
  const VideoDetailsPage(this.model, {super.key});

  final Video model;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          "Video Details",
          style: Theme.of(context).textTheme.titleLarge,
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
          //infoSection(context, "Title", model.title),
          //infoSection(context, "Channel", model.channel),
          //infoSection(context, "Channel ID", model.channelId),
          //infoSection(context, "Upload Date", model.uploadDate),
          //infoSection(context, "Resolution", "${model.width}x${model.height}"),
          //infoSection(context, "FPS", model.fps.toString()),
          //infoSection(context, "Duration", model.duration.toString()),
          infoSection(context, "Blacklisted", model.blacklist.toString()),
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
            style: Theme.of(context).textTheme.titleLarge,
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
