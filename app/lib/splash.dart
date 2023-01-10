import 'dart:async';

import 'package:t4vd/model.dart';
import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';

class SplashPage extends StatefulWidget {
  const SplashPage({super.key});

  @override
  State<SplashPage> createState() => _SplashPageState();
}

class _SplashPageState extends State<SplashPage> {
  bool _phase = true;
  final String _message = 'Loading app...';
  Timer? _timer;

  final delay = const Duration(milliseconds: 1200);

  @override
  void initState() {
    super.initState();
    setTimer();
    WidgetsBinding.instance
        .addPostFrameCallback((timeStamp) async => await ensure());
  }

  Future<void> ensure() async {
    if (!mounted) return;
    final model = ScopedModel.of<BJJModel>(context);
    await model.readCachedBrightness();
    await model.readCachedCreds();
    if (!model.isLoggedIn) {
      if (!mounted) return;
      await Navigator.of(context).pushNamed('/user/login');
    }
    if (!model.hasProject) {
      if (!mounted) return;
      // present the Select Project page, make the back button
      // go to the login screen.
      await model.loadCachedProject();
      if (!model.hasProject) {
        if (!mounted) return;
        await Navigator.of(context).pushNamed('/project/select');
        if (!model.hasProject) {
          // project selection cancelled
          if (!mounted) return;
          WidgetsBinding.instance
              .addPostFrameCallback((timeStamp) async => await ensure());
          return;
        }
      }
    }
    if (!mounted) return;
    Navigator.of(context).popAndPushNamed('/tabs');
  }

  @override
  void dispose() {
    super.dispose();
    _timer?.cancel();
  }

  void setTimer() {
    _timer?.cancel();
    _timer = Timer(delay, () {
      if (!mounted) return;
      setState(() => _phase = !_phase);
      setTimer();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(12.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              AnimatedOpacity(
                duration: delay,
                opacity: _phase ? 0.37 : 0.1,
                child: Container(
                  constraints: const BoxConstraints(
                    maxWidth: 300,
                  ),
                  child: const AspectRatio(
                    aspectRatio: 1,
                    child: Image(
                      image: AssetImage('assets/laselva.png'),
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 20),
              Padding(
                padding: const EdgeInsets.all(16.0),
                child: Text(
                  _message,
                  style: Theme.of(context).textTheme.headlineSmall,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
