import 'package:t4vd/filter/card_stack.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:scoped_model/scoped_model.dart';
import 'package:swipeable_card_stack/swipeable_card_stack.dart';

import '../api.dart';
import '../context.dart';
import '../model.dart';

class CardView extends StatelessWidget {
  const CardView(
    this.marker, {
    super.key,
  });

  final Marker marker;

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Container(
          decoration: BoxDecoration(
              boxShadow: [
                BoxShadow(
                  blurRadius: 12.0,
                  blurStyle: BlurStyle.normal,
                  color: Colors.black.withAlpha(100),
                  offset: const Offset(10, 16),
                )
              ],
              color: Colors.grey,
              border: Border.all(
                color: Theme.of(context).dividerColor,
                width: 1.0,
              )),
          child: Padding(
            padding: const EdgeInsets.all(8.0),
            child: Align(
              child: AspectRatio(
                aspectRatio: 1920.0 / 1080.0,
                child: Container(
                  decoration: BoxDecoration(
                    image: DecorationImage(
                      image: NetworkImage(videoThumbnail(marker.videoId)),
                      alignment: const Alignment(0, 0),
                      fit: BoxFit.cover,
                    ),
                  ),
                ),
              ),
            ),
          ),
        ),
        Positioned(
          top: 0,
          right: 0,
          child: GestureDetector(
            onTapDown: (details) {
              showVideoContextMenu(
                context,
                details.globalPosition,
                id: marker.videoId,
                onBlacklist: () async {},
              );
            },
            child: Icon(
              Icons.more_vert,
              shadows: [
                Shadow(
                  color: Colors.black.withAlpha(170),
                  blurRadius: 12.0,
                  offset: const Offset(0, 0),
                )
              ],
              color: Colors.pink,
            ),
          ),
        )
      ],
    );
  }
}

class FilterPage extends StatefulWidget {
  const FilterPage({super.key});

  @override
  State<FilterPage> createState() => _FilterPageState();
}

class _FilterPageState extends State<FilterPage> {
  final focusNode = FocusNode();
  final SwipeableCardSectionController _cardController =
      SwipeableCardSectionController();
  bool loading = false;

  @override
  void initState() {
    super.initState();
    final model = ScopedModel.of<BJJModel>(context);
    if (model.markers.isEmpty) {
      loading = true;
      model.refreshMarkers(Navigator.of(context)).then((value) {
        if (!mounted) return;
        setState(() => loading = false);
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return ScopedModelDescendant<BJJModel>(
      builder: (context, child, model) {
        return RawKeyboardListener(
          focusNode: focusNode,
          autofocus: true,
          onKey: (ev) {
            if (ev is RawKeyUpEvent) {
              if (ev.logicalKey == LogicalKeyboardKey.keyA) {
                _cardController.triggerSwipeLeft();
              } else if (ev.logicalKey == LogicalKeyboardKey.keyD) {
                _cardController.triggerSwipeRight();
              } else if (ev.logicalKey == LogicalKeyboardKey.keyR) {
                setState(() {});
              }
            }
          },
          child: Stack(
            children: [
              /*
              Container(
                color: Theme.of(context).canvasColor,
                child: Column(
                  children: [
                    SwipeableCardsSection(
                      cardController: _cardController,
                      context: context,
                      items: model.markers
                          .map((marker) => CardView(marker))
                          .toList(),
                      onCardSwiped: (dir, index, widget) {
                        print('swipe $index $dir');
                        //model
                        //    .nextMarker()
                        //    .then((value) => print('got next marker'));
                        //_cardController
                        //    .addItem(CardView(await model.nextMarker()));
                      },
                      enableSwipeUp: false,
                      enableSwipeDown: false,
                    ),
                  ],
                ),
              ),
              */
              const CardStack(),
              Visibility(
                visible: loading,
                child: const Center(
                  child: CircularProgressIndicator(),
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}
