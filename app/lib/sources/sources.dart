import 'package:t4vd/sources/input.dart';
import 'package:t4vd/sources/output.dart';
import 'package:flutter/material.dart';

class SourcesPage extends StatelessWidget {
  const SourcesPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
        color: Colors.green,
        child: const DefaultTabController(
          length: 2,
          child: Scaffold(
            body: TabBarView(children: [
              InputPage(),
              OutputPage(),
            ]),
            bottomNavigationBar: BottomAppBar(
              child: TabBar(
                tabs: [
                  Tab(
                    key: Key('inputTab'),
                    icon: Icon(Icons.input),
                    text: "Input",
                  ),
                  Tab(
                    key: Key('outputTab'),
                    icon: Icon(Icons.output),
                    text: "Output",
                  ),
                ],
              ),
            ),
          ),
        ));
  }
}
