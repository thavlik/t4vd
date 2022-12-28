import 'package:bjjv/account.dart';
import 'package:bjjv/project/create_project.dart';
import 'package:bjjv/project/manage_collaborators.dart';
import 'package:bjjv/project/select_project.dart';
import 'package:bjjv/splash.dart';
import 'package:bjjv/user/login.dart';
import 'package:bjjv/model.dart';
import 'package:bjjv/filter/filter.dart';
import 'package:bjjv/sources/channels.dart';
import 'package:bjjv/sources/playlists.dart';
import 'package:bjjv/sources/videos.dart';
import 'package:bjjv/sources/sources.dart';
import 'package:bjjv/user/preferences.dart';
import 'package:bjjv/user/privacy_policy.dart';
import 'package:bjjv/user/reset_password.dart';
import 'package:bjjv/user/sign_up.dart';
import 'package:flutter/material.dart';
import 'package:scoped_model/scoped_model.dart';

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
        length: 3,
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
                  key: Key('accountTab'),
                  icon: Icon(Icons.account_circle),
                  text: "Account",
                ),
                /*
                Tab(
                  icon: Icon(Icons.rectangle_outlined),
                  text: "Crop",
                ),
                Tab(
                  icon: Icon(Icons.tag),
                  text: "Tags",
                ),
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
              AccountPage(),
              //CropPage(),
              //TagsPage(),
              //LandmarksPage(),
            ],
          ),
        ));
  }
}
