import 'package:t4vd/account.dart';
import 'package:t4vd/graph/graph.dart';
import 'package:t4vd/paint/paint.dart';
import 'package:t4vd/project/create_project.dart';
import 'package:t4vd/project/manage_collaborators.dart';
import 'package:t4vd/project/select_project.dart';
import 'package:t4vd/splash.dart';
import 'package:t4vd/tags.dart';
import 'package:t4vd/user/login.dart';
import 'package:t4vd/model.dart';
import 'package:t4vd/filter/filter.dart';
import 'package:t4vd/sources/channels.dart';
import 'package:t4vd/sources/playlists.dart';
import 'package:t4vd/sources/videos.dart';
import 'package:t4vd/sources/sources.dart';
import 'package:t4vd/user/preferences.dart';
import 'package:t4vd/user/privacy_policy.dart';
import 'package:t4vd/user/reset_password.dart';
import 'package:t4vd/user/sign_up.dart';
import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';

import 'crop.dart';

void main() {
  run();
}

void run({
  String? initialRoute,
  UserCredentials? creds,
}) {
  runApp(MyApp(
    initialRoute: initialRoute,
    initialCreds: creds,
  ));
}

class MyApp extends StatelessWidget {
  final UserCredentials? initialCreds;
  final String? initialRoute;

  const MyApp({
    super.key,
    this.initialCreds,
    this.initialRoute,
  });

  @override
  Widget build(BuildContext context) {
    return ScopedModel<BJJModel>(
      model: BJJModel(creds: initialCreds),
      child: ScopedModelDescendant(
        builder: (BuildContext context, Widget? child, BJJModel model) =>
            MaterialApp(
          initialRoute:
              initialRoute ?? '/splash', //'/project/collaborators', //
          title: 'Brazilian Jiu-jitsu Vision Project',
          themeMode: ThemeMode.dark,
          darkTheme: ThemeData(
            brightness: model.brightness,
            /* dark theme settings */
          ),
          theme: ThemeData(
            primarySwatch: Colors.blue,
          ),
          routes: {
            '/tabs': (context) => const TabsPage(),
            '/project/select': (context) =>
                const SelectProjectPageArgsExtrator(),
            '/project/create': (context) => const CreateProjectPage(),
            '/project/collaborators': (context) =>
                const ManageCollaboratorsPage(),
            '/splash': (context) => const SplashPage(),
            '/user/login': (context) => const LoginPage(),
            '/user/resetpassword': (context) => const ResetPasswordPage(),
            '/user/signup': (context) => const SignUpPage(),
            '/user/privacypolicy': (context) => const PrivacyPolicyPage(),
            '/user/preferences': (context) => const PreferencesPage(),
            '/input/channels': (context) => const InputChannelsPage(),
            '/input/playlists': (context) => const InputPlaylistsPage(),
            '/input/videos': (context) => const InputVideosPage(),
          },
        ),
      ),
    );
  }
}

class TabsPage extends StatefulWidget {
  const TabsPage({super.key});

  @override
  State<TabsPage> createState() => _TabsPageState();
}

class _TabsPageState extends State<TabsPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((timeStamp) async {
      if (!mounted) return;
      final model = ScopedModel.of<BJJModel>(context);
      if (!model.hasProject) {
        await Navigator.of(context).pushNamed('/project/select');
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
        length: 7,
        child: Scaffold(
          appBar: AppBar(
            toolbarHeight: 16,
            automaticallyImplyLeading: false,
            bottom: const TabBar(
              tabs: [
                Tab(
                  key: Key('sourcesTab'),
                  icon: Icon(Icons.source),
                  text: "Sources",
                ),
                Tab(
                  key: Key('filterTab'),
                  icon: Icon(Icons.filter),
                  text: "Filter",
                ),
                Tab(
                  icon: Icon(Icons.tag),
                  text: "Tags",
                ),
                Tab(
                  icon: Icon(Icons.crop),
                  text: "Crop",
                ),
                Tab(
                  icon: Icon(Icons.format_paint),
                  text: "Paint",
                ),
                Tab(
                  icon: Icon(Icons.graphic_eq),
                  text: "Graph",
                ),
                Tab(
                  key: Key('accountTab'),
                  icon: Icon(Icons.account_circle),
                  text: "Account",
                ),
                /*
                Tab(
                  icon: Icon(Icons.edit),
                  text: "Landmarks",
                ),*/
              ],
            ),
          ),
          body: const TabBarView(
            physics: NeverScrollableScrollPhysics(),
            children: [
              SourcesPage(),
              FilterPage(),
              TagsPage(),
              CropPage(),
              PaintPage(),
              GraphPage(),
              AccountPage(),
              //LandmarksPage(),
            ],
          ),
        ));
  }
}
