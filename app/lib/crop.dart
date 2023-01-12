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
  Widget build(BuildContext context) => ScopedModelDescendant<BJJModel>(
        builder: (context, child, model) {
          final currentMarker = model.currentMarker;
          return Stack(
            fit: StackFit.expand,
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.start,
                children: [
                  if (currentMarker != null)
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
                              image: DecorationImage(
                                image: NetworkImage(currentMarker.imageUrl),
                                alignment: const Alignment(0, 0),
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
        },
      );
}
