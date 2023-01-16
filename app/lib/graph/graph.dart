import 'dart:math';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

extension StringExtension on String {
  String capitalize() {
    return "${this[0].toUpperCase()}${substring(1).toLowerCase()}";
  }
}

class NodeDataRef {
  final String id;
  final String channel;

  NodeDataRef({
    required this.id,
    required this.channel,
  });
}

class NodeChannelMeta {
  final String name;
  final Color color;

  NodeChannelMeta({
    required this.name,
    required this.color,
  });
}

class NodeMeta {
  final String name;
  final Color color;
  final List<NodeChannelMeta> inputs;
  final List<NodeChannelMeta> outputs;

  NodeMeta({
    required this.name,
    required this.color,
    this.inputs = const [],
    this.outputs = const [],
  });
}

class Node {
  final NodeMeta meta;
  final String id;
  Map<String, NodeDataRef?> inputLinks;

  void setInputRef(String channel, NodeDataRef ref) {
    if (!meta.inputs.any((e) => e.name == channel)) {
      throw Exception("input channel $channel does not exist");
    }
    inputLinks[channel] = ref;
  }

  Node({
    required this.meta,
    required this.id,
    this.inputLinks = const {},
  });
}

class Graph {
  final List<Node> nodes = [
    Node(
      id: "0",
      meta: NodeMeta(
        name: "Video Dataset",
        color: Colors.red,
        outputs: [
          NodeChannelMeta(
            name: "videos",
            color: Colors.green,
          ),
          NodeChannelMeta(
            name: "frames",
            color: Colors.blue,
          ),
        ],
      ),
    ),
    Node(
      id: "1",
      meta: NodeMeta(
        name: "Filter",
        color: Colors.deepPurple,
        inputs: [
          NodeChannelMeta(
            name: "default",
            color: Colors.green,
          )
        ],
        outputs: [
          NodeChannelMeta(
            name: "default",
            color: Colors.blue,
          ),
          NodeChannelMeta(
            name: "keep",
            color: Colors.green,
          ),
          NodeChannelMeta(
            name: "discard",
            color: Colors.red,
          ),
        ],
      ),
      inputLinks: {
        "default": NodeDataRef(
          id: "0",
          channel: "frames",
        ),
      },
    ),
  ];
}

class _NodeWidget extends StatelessWidget {
  const _NodeWidget(
    this.node, {
    this.selected = false,
    required this.nodeKey,
    required this.nibSize,
    required this.onSelected,
    super.key,
  });

  final NodeKey nodeKey;
  final double nibSize;
  final bool selected;
  final Node node;
  final Function(Node) onSelected;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(96.0),
      child: GestureDetector(
        onTap: () => onSelected(node),
        child: Container(
          key: nodeKey.root,
          width: 200,
          height: 200,
          decoration: BoxDecoration(
            color: Theme.of(context).cardColor,
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.25),
                spreadRadius: 1,
                blurRadius: 4,
                offset: const Offset(0, 4),
              ),
            ],
            border: Border.all(
              color: selected ? Colors.yellow : Colors.black,
              width: 1,
            ),
          ),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Container(
                decoration: BoxDecoration(
                  color: node.meta.color,
                  border: Border(
                    bottom: BorderSide(
                      color: selected ? Colors.yellow : Colors.black,
                      width: 1,
                    ),
                  ),
                ),
                child: Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.center,
                    children: [
                      Text(
                        node.meta.name,
                        style: Theme.of(context).textTheme.bodyLarge,
                      ),
                    ],
                  ),
                ),
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Column(
                    mainAxisAlignment: MainAxisAlignment.start,
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: node.meta.inputs
                        .map((input) => Padding(
                              padding:
                                  const EdgeInsets.symmetric(vertical: 4.0),
                              child: Row(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Container(
                                    key: nodeKey.inputKeys[input.name],
                                    width: nibSize,
                                    height: nibSize,
                                    decoration: BoxDecoration(
                                      color: input.color,
                                      border: Border.all(
                                        color: Colors.black,
                                        width: 1,
                                      ),
                                    ),
                                  ),
                                  const SizedBox(width: 4),
                                  Text(
                                    input.name.capitalize(),
                                    style:
                                        Theme.of(context).textTheme.bodySmall,
                                  ),
                                ],
                              ),
                            ))
                        .toList(),
                  ),
                  Column(
                    mainAxisAlignment: MainAxisAlignment.start,
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: node.meta.outputs
                        .map((output) => Padding(
                              padding:
                                  const EdgeInsets.symmetric(vertical: 4.0),
                              child: Row(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(
                                    output.name.capitalize(),
                                    style:
                                        Theme.of(context).textTheme.bodySmall,
                                  ),
                                  const SizedBox(width: 4),
                                  Container(
                                    key: nodeKey.outputKeys[output.name],
                                    width: nibSize,
                                    height: nibSize,
                                    decoration: BoxDecoration(
                                      color: output.color,
                                      border: Border.all(
                                        color: Colors.black,
                                        width: 1,
                                      ),
                                    ),
                                  ),
                                ],
                              ),
                            ))
                        .toList(),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class GraphPage extends StatefulWidget {
  const GraphPage({super.key});

  @override
  State<GraphPage> createState() => _GraphPageState();
}

class _Candidate {
  Offset location = Offset.zero;
}

class NodeKey {
  GlobalKey root = GlobalKey();
  Map<String, GlobalKey> inputKeys = {};
  Map<String, GlobalKey> outputKeys = {};
}

class NodeKeys {
  final Map<String, NodeKey> nodes = {};
}

class _GraphPageState extends State<GraphPage> {
  final focusNode = FocusNode();
  final origin = GlobalKey();
  final graph = Graph();
  Set<String> selection = {"0"};
  bool additive = false;
  bool group = false;
  final nibSize = 16.0;
  final NodeKeys keys = NodeKeys();

  @override
  void initState() {
    super.initState();
    focusNode.requestFocus();
  }

  void onKeyEvent(KeyEvent event) {
    if (event is KeyDownEvent) {
      if (event.logicalKey == LogicalKeyboardKey.shiftLeft) {
        additive = true;
      } else if (event.logicalKey == LogicalKeyboardKey.controlLeft) {
        group = true;
      }
      return;
    }
    if (event is KeyUpEvent) {
      if (event.logicalKey == LogicalKeyboardKey.shiftLeft) {
        additive = false;
      } else if (event.logicalKey == LogicalKeyboardKey.controlLeft) {
        group = false;
      }
      return;
    }
  }

  void select(Node node) => setState(() {
        if (additive) {
          if (selection.contains(node.id)) {
            selection.remove(node.id);
            return;
          }
          selection.add(node.id);
        } else {
          selection = {node.id};
        }
      });

  RenderBox get originRenderBox =>
      origin.currentContext!.findRenderObject() as RenderBox;

  Offset get originOffset =>
      originRenderBox.localToGlobal(Offset.zero).scale(-1, -1);

  _HitTestResult? _source;

  void onPointerDown(PointerDownEvent ev) {
    _source = _hitTest(
      graph: graph,
      nibSize: nibSize,
      originOffset: originOffset,
      testPosition: ev.localPosition,
      keys: keys,
    );
    if (_source != null) {
      // TODO: set candidate connection
      return;
    }
  }

  void onPointerUp(PointerUpEvent ev) {
    final sink = _hitTest(
      graph: graph,
      nibSize: nibSize,
      originOffset: originOffset,
      testPosition: ev.localPosition,
      keys: keys,
    );
    if (sink != null && _source != null) {
      connect(_source!, sink);
    }
    _source = null;
    // TODO: clear candidate connections
  }

  void connect(_HitTestResult source, _HitTestResult sink) {
    if (source.node == sink.node ||
        source.channel == null ||
        sink.channel == null) return;
    if (!source.isOutput) {
      final temp = source;
      source = sink;
      sink = temp;
    }
    if (!source.isOutput || sink.isOutput) return;
    // TODO: try and make connection in backend
    setState(() => sink.node.inputLinks[sink.channel!] = NodeDataRef(
          id: source.node.id,
          channel: source.channel!,
        ));
  }

  void onPointerMove(PointerMoveEvent ev) {
    if (_source != null) {
      // TODO: update candidate connection position
    }
  }

  void buildKeys() {
    keys.nodes.clear();
    for (final node in graph.nodes) {
      final n = NodeKey();
      for (final input in node.meta.inputs) {
        n.inputKeys[input.name] = GlobalKey();
      }
      for (final output in node.meta.outputs) {
        n.outputKeys[output.name] = GlobalKey();
      }
      keys.nodes[node.id] = n;
    }
  }

  @override
  Widget build(BuildContext context) {
    buildKeys();
    return Listener(
      onPointerDown: onPointerDown,
      onPointerUp: onPointerUp,
      onPointerMove: onPointerMove,
      child: GestureDetector(
        onTap: () => setState(() => selection = {}),
        child: KeyboardListener(
          focusNode: focusNode,
          onKeyEvent: onKeyEvent,
          child: Scaffold(
            key: origin,
            body: Stack(
              children: [
                CustomPaint(
                  painter: MyFancyPainter(
                    origin: origin,
                    nibSize: nibSize,
                    graph: graph,
                    nodeKeys: keys,
                  ),
                ),
                Row(
                  children: graph.nodes
                      .map((node) => _NodeWidget(
                            node,
                            nibSize: nibSize,
                            selected: selection.contains(node.id),
                            onSelected: select,
                            nodeKey: keys.nodes[node.id]!,
                          ))
                      .toList(),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

class _HitTestResult {
  final Node node;
  final bool isOutput;
  final String? channel;

  _HitTestResult({
    required this.node,
    this.isOutput = false,
    this.channel,
  });

  @override
  String toString() => channel == null
      ? 'HitTestResult(${node.id}, ${node.meta.name})'
      : 'HitTestResult(${node.id}, ${node.meta.name}, ${isOutput ? "output" : "input"}:$channel)';
}

_HitTestResult? _hitTest({
  required Graph graph,
  required Offset originOffset,
  required Offset testPosition,
  required double nibSize,
  required NodeKeys keys,
}) {
  for (var entry in keys.nodes.entries) {
    final nodeId = entry.key;
    final node = graph.nodes.firstWhere((node) => node.id == nodeId);
    final nodeKey = entry.value;
    for (var inputEntry in nodeKey.inputKeys.entries) {
      final inputName = inputEntry.key;
      final inputKey = inputEntry.value;
      final inputRenderBox =
          inputKey.currentContext!.findRenderObject() as RenderBox;
      final inputOffset = inputRenderBox.localToGlobal(originOffset);
      final inputRect = Rect.fromLTWH(
        inputOffset.dx,
        inputOffset.dy,
        nibSize,
        nibSize,
      );
      if (inputRect.contains(testPosition)) {
        return _HitTestResult(
          node: node,
          isOutput: false,
          channel: inputName,
        );
      }
    }
    for (var outputEntry in nodeKey.outputKeys.entries) {
      final outputName = outputEntry.key;
      final outputKey = outputEntry.value;
      final outputRenderBox =
          outputKey.currentContext!.findRenderObject() as RenderBox;
      final outputOffset = outputRenderBox.localToGlobal(originOffset);
      final outputRect = Rect.fromLTWH(
        outputOffset.dx,
        outputOffset.dy,
        nibSize,
        nibSize,
      );
      if (outputRect.contains(testPosition)) {
        return _HitTestResult(
          node: node,
          isOutput: true,
          channel: outputName,
        );
      }
    }
    final rootRenderBox =
        nodeKey.root.currentContext!.findRenderObject() as RenderBox;
    final rootOffset = rootRenderBox.localToGlobal(originOffset);
    final rootRect = Rect.fromLTWH(
      rootOffset.dx,
      rootOffset.dy,
      rootRenderBox.size.width,
      rootRenderBox.size.height,
    );
    if (rootRect.contains(testPosition)) {
      return _HitTestResult(node: node);
    }
  }
  return null;
}

class MyFancyPainter extends CustomPainter {
  final GlobalKey origin;
  final double nibSize;
  final Graph graph;
  final NodeKeys nodeKeys;

  MyFancyPainter({
    required this.origin,
    required this.nibSize,
    required this.graph,
    required this.nodeKeys,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..strokeWidth = 2
      ..color = Colors.white.withAlpha(128)
      ..style = PaintingStyle.stroke;
    //final arc = Path();
    //arc.moveTo(200, 200);
    //arc.cubicTo(400, 200, 400, 400, 600, 400);
    //canvas.drawPath(arc, paint);
    final originRenderBox =
        origin.currentContext!.findRenderObject() as RenderBox;
    final originOffset =
        originRenderBox.localToGlobal(Offset.zero).scale(-1, -1);
    final halfNibSize = nibSize / 2;
    graph.nodes.where((node) => node.inputLinks.isNotEmpty).forEach((node) {
      node.inputLinks.forEach((key, value) {
        if (value == null) return;
        final sourceKey = nodeKeys.nodes[node.id]!;
        final sinkKey = nodeKeys.nodes[value.id]!;
        final inputKey = sourceKey.inputKeys[key]!;
        final outputKey = sinkKey.outputKeys[value.channel]!;
        final inputRenderBox =
            inputKey.currentContext!.findRenderObject() as RenderBox;
        final outputRenderBox =
            outputKey.currentContext!.findRenderObject() as RenderBox;
        final inputOffset = inputRenderBox.localToGlobal(originOffset);
        final outputOffset = outputRenderBox.localToGlobal(originOffset);
        final path = Path();
        path.moveTo(inputOffset.dx + halfNibSize, inputOffset.dy + halfNibSize);
        final hx = (inputOffset.dx + outputOffset.dx) / 2 + halfNibSize;
        final hy = (inputOffset.dy + outputOffset.dy) / 2 + halfNibSize;
        path.cubicTo(
          hx + halfNibSize,
          inputOffset.dy + halfNibSize,
          hx,
          hy,
          outputOffset.dx + halfNibSize,
          outputOffset.dy + halfNibSize,
        );
        canvas.drawPath(path, paint);
      });
    });
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => true;
}
