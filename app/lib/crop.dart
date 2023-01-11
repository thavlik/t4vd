import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';

import 'model.dart';

class CropPage extends StatelessWidget {
  const CropPage({super.key});

  Future<void> submit(BuildContext context) async {}

  Future<void> discard(BuildContext context) async =>
      await ScopedModel.of<BJJModel>(context).discard(Navigator.of(context));

  void previous(BuildContext context) =>
      ScopedModel.of<BJJModel>(context).markerBack();

  @override
  Widget build(BuildContext context) {
    return Stack(
      fit: StackFit.expand,
      children: [
        Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            Expanded(
              child: Align(
                child: AspectRatio(
                  aspectRatio: 1920.0 / 1080.0,
                  child: Container(
                    decoration: BoxDecoration(
                      border: Border.all(
                        color: Colors.black.withAlpha(180),
                        width: 32.0,
                      ),
                      image: const DecorationImage(
                        image: AssetImage("assets/example-0.jpg"),
                        alignment: Alignment(0, 0),
                        fit: BoxFit.cover,
                      ),
                    ),
                  ),
                ),
              ),
            ),
          ],
        ),
        Visibility(
          visible: ScopedModel.of<BJJModel>(context).markerIndex > 0,
          child: Positioned(
              top: 16,
              left: 16,
              child: FloatingActionButton(
                onPressed: () {},
                child: const Icon(Icons.navigate_before),
              )),
        ),
        Positioned(
            bottom: 16,
            left: 16,
            child: FloatingActionButton(
              onPressed: () {},
              child: const Icon(Icons.block),
            )),
        Positioned(
            bottom: 16,
            right: 16,
            child: FloatingActionButton(
              onPressed: () {},
              child: const Icon(Icons.done),
            )),
      ],
    );
  }
}
