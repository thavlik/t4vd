import 'package:flutter/material.dart';

import '../api.dart';

TextStyle headerStyle(BuildContext context) {
  return Theme.of(context)
      .textTheme
      .titleSmall!
      .copyWith(fontWeight: FontWeight.bold);
}

class SectionHeader extends StatelessWidget {
  const SectionHeader({super.key, required this.text});

  final String text;

  @override
  Widget build(BuildContext context) {
    return Text(text, style: Theme.of(context).textTheme.headline5);
  }
}

class ChannelsSection extends StatelessWidget {
  const ChannelsSection({super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(8.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const SectionHeader(text: "Channels"),
          Table(
            children: [
              TableRow(
                children: [
                  TableCell(child: Text("ID", style: headerStyle(context))),
                  TableCell(child: Text("Name", style: headerStyle(context))),
                  TableCell(
                      child: Text("# Videos", style: headerStyle(context))),
                  TableCell(
                      child:
                          Text("Last Retrieved", style: headerStyle(context))),
                ],
              ),
            ],
          )
        ],
      ),
    );
  }
}

class PlaylistsSection extends StatelessWidget {
  const PlaylistsSection({super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(8.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const SectionHeader(text: "Playlists"),
          Table(
            children: [
              TableRow(
                children: [
                  TableCell(child: Text("ID", style: headerStyle(context))),
                  TableCell(child: Text("Title", style: headerStyle(context))),
                  TableCell(
                      child: Text("Channel", style: headerStyle(context))),
                  TableCell(
                      child: Text("# Videos", style: headerStyle(context))),
                  TableCell(
                      child:
                          Text("Last Retrieved", style: headerStyle(context))),
                ],
              ),
            ],
          )
        ],
      ),
    );
  }
}

class VideosSection extends StatelessWidget {
  const VideosSection(this.videos, {super.key});

  final List<VideoListItem> videos;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(8.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const SectionHeader(text: "Videos"),
          Table(
            children: [
              TableRow(
                children: [
                  TableCell(child: Text("ID", style: headerStyle(context))),
                  TableCell(child: Text("Title", style: headerStyle(context))),
                  TableCell(
                      child: Text("Channel", style: headerStyle(context))),
                  TableCell(
                      child: Text("Duration", style: headerStyle(context))),
                  //TableCell(child: Text("Uploaded", style: headerStyle(context))),
                  TableCell(child: Text("Width", style: headerStyle(context))),
                  TableCell(child: Text("Height", style: headerStyle(context))),
                  TableCell(child: Text("FPS", style: headerStyle(context))),
                ],
              ),
              ...videos
                  .map((video) => TableRow(
                        children: [
                          TableCell(child: Text(video.id)),
                          TableCell(
                            child: Text(
                              video.title,
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                          TableCell(
                            child: Text(
                              video.channel,
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                          TableCell(child: Text(video.duration.toString())),
                          //TableCell(child: Text(video.uploadDate)),
                          TableCell(child: Text(video.width.toString())),
                          TableCell(child: Text(video.height.toString())),
                          TableCell(child: Text(video.fps.toString())),
                        ],
                      ))
                  .toList(),
            ],
          ),
        ],
      ),
    );
  }
}

class InputPage extends StatelessWidget {
  const InputPage({super.key});

  Widget listItem({
    required BuildContext context,
    required String name,
    required void Function() onTap,
    Key? key,
  }) {
    return InkWell(
      key: key,
      onTap: onTap,
      child: Container(
        decoration: BoxDecoration(
          border: Border(
            bottom: BorderSide(color: Theme.of(context).dividerColor),
          ),
        ),
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Padding(
                padding: const EdgeInsets.all(16.0),
                child: Text(
                  name,
                  style: Theme.of(context).textTheme.headline4,
                ),
              ),
              const Padding(
                padding: EdgeInsets.all(8.0),
                child: Icon(Icons.navigate_next),
              ),
            ],
          ),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        listItem(
          context: context,
          name: "Channels",
          onTap: () => Navigator.of(context).pushNamed("/input/channels"),
          key: const Key('channels'),
        ),
        listItem(
          context: context,
          name: "Playlists",
          onTap: () => Navigator.of(context).pushNamed("/input/playlists"),
          key: const Key('playlists'),
        ),
        listItem(
          context: context,
          name: "Videos",
          onTap: () => Navigator.of(context).pushNamed("/input/videos"),
          key: const Key('videos'),
        ),
      ],
    );
  }
}
