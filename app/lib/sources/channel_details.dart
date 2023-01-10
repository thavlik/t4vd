import 'package:flutter/material.dart';

import '../api.dart';

class ChannelDetailsPage extends StatelessWidget {
  const ChannelDetailsPage(this.model, {super.key});

  final Channel model;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          "Channel Details",
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
          infoSection(
              context, "Name", model.info?.name ?? '<currently unavailable>'),
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
