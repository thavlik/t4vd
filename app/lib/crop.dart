import 'package:flutter/material.dart';

class CropPage extends StatelessWidget {
  const CropPage({super.key});

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
        Positioned(
            bottom: 16,
            left: 16,
            child: FloatingActionButton(
              onPressed: () {},
              child: const Icon(Icons.navigate_before),
            )),
        Positioned(
            top: 16,
            left: 16,
            child: FloatingActionButton(
              onPressed: () {},
              child: const Icon(Icons.cancel_outlined),
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
