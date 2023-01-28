import 'package:flutter/widgets.dart';
import 'package:scoped_model/scoped_model.dart';

import '../api.dart';
import '../model.dart';

class PaintPage extends StatefulWidget {
  const PaintPage({super.key});

  @override
  State<PaintPage> createState() => _PaintPageState();
}

class _PaintPageState extends State<PaintPage> {
  bool _loading = false;

  @override
  void initState() {
    super.initState();
    final model = ScopedModel.of<BJJModel>(context);
    WidgetsBinding.instance.addPostFrameCallback((timeStamp) async {
      if (!mounted) return;
      if (model.markers == null ||
          model.markers!.isEmpty ||
          model.markerIndex == model.markers!.length) {
        setState(() => _loading = true);
        try {
          await model.refreshMarkers(Navigator.of(context));
        } on InvalidCredentialsError {
          Navigator.of(context).pushNamed('/splash');
        } on ForbiddenError {
          Navigator.of(context).pushNamed('/splash');
        } finally {}
        if (!mounted) return;
        setState(() => _loading = false);
      }
      //if (!mounted) return;
      //model.precacheFrames(context);
    });
  }

  @override
  Widget build(BuildContext context) {
    return const Placeholder();
  }
}
