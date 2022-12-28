import 'package:t4vd/sources/videos.dart';
import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';

import '../model.dart';

class OutputPage extends StatefulWidget {
  const OutputPage({super.key});

  @override
  State<OutputPage> createState() => _OutputPageState();
}

class _OutputPageState extends State<OutputPage> {
  bool loading = true;

  @override
  void initState() {
    super.initState();
    ScopedModel.of<BJJModel>(context)
        .refreshDataset(context)
        .then((value) => setState(() => loading = false));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: ScopedModelDescendant<BJJModel>(builder: (context, child, model) {
        return Stack(
          children: [
            model.dataset == null
                ? Visibility(
                    visible: !loading,
                    child: Padding(
                      padding: const EdgeInsets.all(16.0),
                      child: Opacity(
                        opacity: 0.5,
                        child: Text(
                          "There are no videos in the output dataset.",
                          style: Theme.of(context).textTheme.bodyMedium,
                        ),
                      ),
                    ),
                  )
                : ListView(
                    children: [
                      ...model.dataset!.videos
                          .map((vid) => InputVideoListItem(
                                model: vid,
                              ))
                          .toList(),
                    ],
                  ),
            Visibility(
              visible: loading,
              child: const Center(
                child: CircularProgressIndicator(),
              ),
            ),
          ],
        );
      }),
    );
  }
}
