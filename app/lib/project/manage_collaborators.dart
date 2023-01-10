import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';

import '../api.dart';
import '../model.dart';

class ManageCollaboratorsPage extends StatefulWidget {
  const ManageCollaboratorsPage({super.key});

  @override
  State<ManageCollaboratorsPage> createState() =>
      _ManageCollaboratorsPageState();
}

class _ManageCollaboratorsPageState extends State<ManageCollaboratorsPage> {
  bool _editMode = false;
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((timeStamp) async {
      if (!mounted) return;
      final model = ScopedModel.of<BJJModel>(context);
      await model.readCachedCreds();
      if (!mounted) return;
      await model.ensureProject(Navigator.of(context));
      if (!mounted) return;
      setState(() => _loading = false);
    });
  }

  void addCollaborator(BuildContext context, SearchUser user) async =>
      await ScopedModel.of<BJJModel>(context)
          .addCollaborator(Navigator.of(context), user);

  void removeCollaborator(BuildContext context, SearchUser user) async =>
      await ScopedModel.of<BJJModel>(context)
          .removeCollaborator(Navigator.of(context), user.id);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Manage Collaborators'),
        actions: [
          Tooltip(
            message: 'Remove collaborators',
            child: IconButton(
              onPressed: () => setState(() => _editMode = !_editMode),
              icon: Icon(
                _editMode ? Icons.cancel : Icons.edit,
              ),
            ),
          ),
        ],
      ),
      body: ScopedModelDescendant<BJJModel>(
          builder: (context, child, model) => Center(
                child: Container(
                  constraints: const BoxConstraints(
                    maxWidth: 280,
                  ),
                  child: SingleChildScrollView(
                    child: Padding(
                      padding: const EdgeInsets.all(8.0),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.stretch,
                        children: [
                          Container(
                            decoration: BoxDecoration(
                                border: Border.all(
                              width: 1.0,
                              color: Theme.of(context).dividerColor,
                            )),
                            child: Padding(
                              padding: const EdgeInsets.symmetric(
                                horizontal: 2,
                              ),
                              child: Autocomplete(
                                displayStringForOption: (SearchUser user) =>
                                    user.username,
                                onSelected: (value) =>
                                    addCollaborator(context, value),
                                optionsBuilder:
                                    (TextEditingValue textEditingValue) async =>
                                        textEditingValue.text.length <
                                                minUsernameLength
                                            ? <SearchUser>[]
                                            : await ScopedModel.of<BJJModel>(
                                                    context)
                                                .searchUsers(
                                                    Navigator.of(context),
                                                    textEditingValue.text),
                              ),
                            ),
                          ),
                          ...model.project?.collaborators
                                  ?.map((user) => CollaboratorListItem(
                                        user: user,
                                        editMode: _editMode,
                                        onDelete: () =>
                                            removeCollaborator(context, user),
                                      ))
                                  .toList() ??
                              [],
                          if (_loading ||
                              (model.project?.collaborators?.isEmpty ?? false))
                            Padding(
                              padding: const EdgeInsets.all(4.0),
                              child: Opacity(
                                opacity: 0.7,
                                child: Container(
                                  height: 64,
                                  decoration: BoxDecoration(
                                      border: Border.all(
                                    width: 1,
                                    color: Theme.of(context).dividerColor,
                                  )),
                                  child: _loading
                                      ? const Padding(
                                          padding: EdgeInsets.all(12.0),
                                          child: Center(
                                            child: CircularProgressIndicator(),
                                          ),
                                        )
                                      : model.project?.collaborators?.isEmpty ??
                                              false
                                          ? Padding(
                                              padding:
                                                  const EdgeInsets.all(20.0),
                                              child: Text(
                                                'Project has no collaborators. Add one!',
                                                style: Theme.of(context)
                                                    .textTheme
                                                    .bodySmall,
                                              ),
                                            )
                                          : Container(),
                                ),
                              ),
                            ),
                        ],
                      ),
                    ),
                  ),
                ),
              )),
    );
  }
}

class CollaboratorListItem extends StatelessWidget {
  const CollaboratorListItem({
    super.key,
    required this.user,
    required this.editMode,
    required this.onDelete,
  });

  final SearchUser user;
  final bool editMode;
  final void Function() onDelete;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(4.0),
      child: Material(
        color: Colors.transparent,
        child: Container(
          height: 64,
          decoration: BoxDecoration(
            color: Theme.of(context).colorScheme.background.withAlpha(64),
            border: Border.all(
              width: 1.0,
              color: Theme.of(context).dividerColor,
            ),
          ),
          child: Padding(
            padding: const EdgeInsets.fromLTRB(8, 0, 0, 0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                Text(
                  user.username,
                  style: Theme.of(context).textTheme.bodyMedium,
                ),
                AnimatedOpacity(
                  duration: const Duration(milliseconds: 200),
                  opacity: editMode ? 1.0 : 0.0,
                  child: IconButton(
                    onPressed: editMode ? onDelete : null,
                    icon: const Icon(Icons.delete),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
