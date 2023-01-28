import 'package:flutter/widgets.dart';

import 'api.dart' as api;
import 'model.dart';

class MarkerStack {
  int get markerIndex => _markerIndex;
  List<api.Marker>? get markers => _markers;
  api.Marker? get currentMarker => _markers == null
      ? null
      : _markers!.isNotEmpty && _markerIndex < _markers!.length
          ? _markers![_markerIndex]
          : null;

  final BJJModel _model;
  final int cacheSize;
  int _markerIndex = 0;
  List<api.Marker>? _markers;

  MarkerStack(
    this._model, {
    this.cacheSize = 1,
  });

  // Advances the FrameStack to the next marker.
  // TODO: implement cacheSize > 1
  Future<void> advance(NavigatorState nav) async =>
      await _model.withAuth(nav, () async {
        if (_markers?.isEmpty ?? true) {
          throw StateError('No markers');
        }
        await _model.ensureProject(nav);
        _markerIndex++;
        if (_markers!.length - _markerIndex < cacheSize) {
          // getStack will throw ResourceNotFoundError
          // if the project is empty. The caller is
          // responsible for handling this and showing
          // an error message.
          final stack = await api.getStack(
            projectId: _model.project!.id,
            creds: _model.creds!,
            size: cacheSize,
          );
          _markers ??= [];
          _markers!.addAll(stack);
        }
      });

  // Initializes the state of the FrameStack
  // by fetching the first marker. If the project
  // has no videos, the markers list will be empty.
  Future<void> reset(NavigatorState nav) async =>
      await _model.withAuth(nav, () async {
        await _model.ensureProject(nav);
        try {
          _markers = [
            await api.getRandomMarker(
              projectId: _model.project!.id,
              creds: _model.creds!,
            )
          ];
        } on api.ResourceNotFoundError {
          // the project is empty
          _markers = [];
        }
        _markerIndex = 0;
      });
}
