import 'package:t4vd/project/select_project.dart';
import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';

import 'model.dart';

class AccountPage extends StatefulWidget {
  const AccountPage({super.key});

  @override
  State<AccountPage> createState() => _AccountPageState();
}

class _AccountPageState extends State<AccountPage> {
  void _signOut(BuildContext context) async {
    await ScopedModel.of<BJJModel>(context).signOut();
    if (!mounted) return;
    Navigator.of(context).popAndPushNamed('/splash');
  }

  Future<void> _switchProject(BuildContext context) async =>
      await Navigator.of(context).pushNamed('/project/select',
          arguments: SelectProjectPageArgs(
            navigatorBackBehavior: NavigatorBackBehavior.pop,
          ));

  Future<void> _manageCollaborators(BuildContext context) async =>
      await Navigator.of(context).pushNamed('/project/collaborators');

  Future<void> _preferences(BuildContext context) async =>
      await Navigator.of(context).pushNamed('/user/preferences');

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Container(
          constraints: const BoxConstraints(maxWidth: 300),
          child: SingleChildScrollView(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.start,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Material(
                    color:
                        Theme.of(context).buttonTheme.colorScheme!.background,
                    child: InkWell(
                      key: const Key('manageCollaborators'),
                      onTap: () => _manageCollaborators(context),
                      child: Padding(
                        padding: const EdgeInsets.all(16.0),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: const [
                            Text('Manage collaborators'),
                          ],
                        ),
                      ),
                    ),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Material(
                    color:
                        Theme.of(context).buttonTheme.colorScheme!.background,
                    child: InkWell(
                      key: const Key('switchProject'),
                      onTap: () => _switchProject(context),
                      child: Padding(
                        padding: const EdgeInsets.all(16.0),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: const [
                            Text('Switch project'),
                          ],
                        ),
                      ),
                    ),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Material(
                    color:
                        Theme.of(context).buttonTheme.colorScheme!.background,
                    child: InkWell(
                      key: const Key('preferences'),
                      onTap: () => _preferences(context),
                      child: Padding(
                        padding: const EdgeInsets.all(16.0),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: const [
                            Text('Preferences'),
                          ],
                        ),
                      ),
                    ),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Material(
                    color:
                        Theme.of(context).buttonTheme.colorScheme!.background,
                    child: InkWell(
                      key: const Key('signOut'),
                      onTap: () => _signOut(context),
                      child: Padding(
                        padding: const EdgeInsets.all(16.0),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: const [
                            Text('Sign out'),
                          ],
                        ),
                      ),
                    ),
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
