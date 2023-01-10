import 'dart:async';

import 'package:t4vd/api.dart';
import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';
import '../model.dart';

class CreateProjectPage extends StatefulWidget {
  const CreateProjectPage({super.key});

  @override
  State<CreateProjectPage> createState() => _CreateProjectPageState();
}

class _CreateProjectPageState extends State<CreateProjectPage> {
  static const projectAlreadyExistsMessage =
      'A project with this name already exists';

  final _nameController = TextEditingController();
  final _nameFocus = FocusNode();
  final Map<String, bool> _takenNames = {};
  Future<bool>? _checking;
  String? _nameError;
  bool get isNameValid =>
      _nameError == null &&
      _nameController.text.length >= minProjectNameLength &&
      _checking == null &&
      _takenNames[_nameController.text] == false;
  Timer? _checkTakenName;

  Future<bool?> nameTaken() async {
    if (!mounted) {
      return true;
    }
    _checkTakenName?.cancel();
    final value = _nameController.text;
    bool? taken;
    try {
      taken = await projectExists(value);
      _takenNames[value] = taken;
      setState(() => _nameError = nameError(value));
    } catch (err) {
      setState(() => _nameError = 'Error: ${err.toString()}');
    }
    return taken;
  }

  String? nameError(String value) {
    if (value.length < minProjectNameLength) {
      return 'Must contain at least $minProjectNameLength characters';
    } else if (' \t\n!@#\$%^&*()\'".,?;:<>[]{}/\\'
        .split('')
        .any((c) => value.contains(c))) {
      return 'Cannot contain spaces or special characters';
    } else if (_takenNames[value] == true) {
      return projectAlreadyExistsMessage;
    }
    return null;
  }

  void onNameChanged(String value) async {
    setState(() => _nameError = nameError(value));
    _checkTakenName?.cancel();
    if (value.length >= minProjectNameLength &&
        !_takenNames.containsKey(value)) {
      _checkTakenName = Timer(const Duration(seconds: 1), nameTaken);
    }
  }

  @override
  void dispose() {
    super.dispose();
    _checkTakenName?.cancel();
  }

  void submit(BuildContext context) async {
    if (!isNameValid) return;
    _checkTakenName?.cancel();
    if (await nameTaken() != false) {
      return;
    }
    if (!mounted) return;
    await ScopedModel.of<BJJModel>(context).createProject(
      nav: Navigator.of(context),
      name: _nameController.text,
    );
    if (!mounted) return;
    Navigator.of(context).pop();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Create project'),
        leading: IconButton(
          key: const Key('createProjectNavBack'),
          onPressed: () => Navigator.of(context).pop(),
          icon: const Icon(
            Icons.navigate_before,
          ),
        ),
      ),
      body: Center(
        child: SingleChildScrollView(
            child: Container(
          constraints: const BoxConstraints(maxWidth: 300),
          child: Column(
            children: [
              Padding(
                padding: const EdgeInsets.all(8.0),
                child: TextField(
                  key: const Key('createProjectName'),
                  enableSuggestions: false,
                  autocorrect: false,
                  controller: _nameController,
                  focusNode: _nameFocus,
                  textInputAction: TextInputAction.next,
                  autofocus: true,
                  onChanged: onNameChanged,
                  onSubmitted: (value) {
                    onNameChanged(value);
                    if (!isNameValid) {
                      _nameFocus.requestFocus();
                      return;
                    }
                    submit(context);
                  },
                  decoration: InputDecoration(
                    border: const OutlineInputBorder(),
                    hintText: 'Name',
                    labelText: _nameError ??
                        (_takenNames[_nameController.text] == false
                            ? 'Project name is available'
                            : (_nameController.text.isEmpty
                                ? 'Enter a project name'
                                : 'Checking availability...')),
                    labelStyle: TextStyle(
                      color: _nameError == null
                          ? (_takenNames[_nameController.text] == false
                              ? Colors.green
                              : Colors.white)
                          : Colors.red,
                    ),
                  ),
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(4.0),
                child: Material(
                  color: Theme.of(context).colorScheme.background,
                  child: InkWell(
                    key: const Key('createProjectSubmit'),
                    onTap: isNameValid ? () => submit(context) : null,
                    child: SizedBox(
                      height: 64,
                      child: Padding(
                        padding: const EdgeInsets.all(16.0),
                        child: Center(
                          child: Opacity(
                            opacity: isNameValid ? 1.0 : 0.5,
                            child: Text(
                              !isNameValid
                                  ? 'Create project'
                                  : 'Create project "${_nameController.text}"',
                              style: Theme.of(context).textTheme.bodyLarge,
                            ),
                          ),
                        ),
                      ),
                    ),
                  ),
                ),
              ),
            ],
          ),
        )),
      ),
    );
  }
}
