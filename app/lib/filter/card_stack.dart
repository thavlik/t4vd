import 'package:t4vd/api.dart';
import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';
import 'package:url_launcher/url_launcher.dart';

import '../model.dart';

enum Direction {
  left,
  right,
  up,
  down,
}

class CardStack extends StatefulWidget {
  const CardStack({super.key});

  @override
  State<CardStack> createState() => _CardStackState();
}

class _CardStackState extends State<CardStack> {
  bool _loading = false;

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
        } on InvalidCredentialsError catch (_) {
          Navigator.of(context).pushNamed('/splash');
        } finally {
          setState(() => _loading = false);
        }
      }
      if (!mounted) return;
      //model.precacheFrames(context);
    });
  }

  void submitLabel(BuildContext context, bool label) async {
    final model = ScopedModel.of<BJJModel>(context);
    await model.classify(
      nav: Navigator.of(context),
      label: label,
    );
    if (!mounted) return;
    //model.precacheFrames(context);
  }

  void onPanUpdate(DragUpdateDetails details) {}

  void onPanEnd(DragEndDetails details) {}

  @override
  Widget build(BuildContext context) {
    return ScopedModelDescendant<BJJModel>(builder: (context, child, model) {
      final marker = model.currentMarker;
      return GestureDetector(
        onPanUpdate: onPanUpdate,
        onPanEnd: onPanEnd,
        child: Stack(
          children: [
            if (marker != null)
              Stack(
                children: [
                  Align(
                    child: AspectRatio(
                      aspectRatio: 1920.0 / 1080.0,
                      child: Container(
                        decoration: BoxDecoration(
                          image: DecorationImage(
                            image: NetworkImage(marker.imageUrl),
                            alignment: const Alignment(0, 0),
                            fit: BoxFit.cover,
                          ),
                        ),
                      ),
                    ),
                  ),
                  Positioned(
                      right: 16,
                      top: 16,
                      child: FloatingActionButton(
                        onPressed: () async {
                          final marker = model.currentMarker!;
                          final t = Duration(microseconds: marker.time ~/ 1000);
                          final url = Uri.parse(
                              "https://youtube.com/watch?v=${marker.videoId}&t=${t.inSeconds}s"); // "https://youtube.com/watch?v=${}");
                          if (!await launchUrl(url)) {
                            throw 'Could not launch $url';
                          }
                        },
                        child: const Icon(Icons.open_in_browser),
                      )),
                ],
              ),
            Visibility(
              visible: model.markerIndex > 0,
              child: Positioned(
                  top: 16,
                  left: 16,
                  child: FloatingActionButton(
                    onPressed: () {
                      model.classifyBack();
                    },
                    child: const Icon(Icons.navigate_before),
                  )),
            ),
            Positioned(
                left: 16,
                bottom: 16,
                child: FloatingActionButton(
                  onPressed: () => submitLabel(context, false),
                  child: const Icon(Icons.block),
                )),
            Positioned(
                right: 16,
                bottom: 16,
                child: FloatingActionButton(
                  onPressed: () => submitLabel(context, true),
                  child: const Icon(Icons.done),
                )),
            Visibility(
              visible: _loading,
              child: const Center(
                child: CircularProgressIndicator(),
              ),
            ),
          ],
        ),
      );
    });
  }
}
