import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';

import '../model.dart';

class PreferencesPage extends StatefulWidget {
  const PreferencesPage({super.key});

  @override
  State<PreferencesPage> createState() => _PreferencesPageState();
}

class _PreferencesPageState extends State<PreferencesPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Preferences'),
      ),
      body: ScopedModelDescendant<BJJModel>(
        builder: (context, child, model) => Center(
          child: Container(
            constraints: const BoxConstraints(
              maxWidth: 300,
            ),
            child: SingleChildScrollView(
              child: Column(
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Dark Mode'),
                      Checkbox(
                        value: model.brightness == Brightness.dark,
                        onChanged: (value) async => await model.setBrightness(
                            (value ?? false)
                                ? Brightness.dark
                                : Brightness.light),
                      )
                    ],
                  )
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
