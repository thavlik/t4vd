import 'dart:async';

import 'package:t4vd/api.dart';
import 'package:t4vd/model.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_test/flutter_test.dart';
//import 'package:flutter_driver/flutter_driver.dart' as fd;
import 'package:integration_test/integration_test.dart';
import 'package:t4vd/main.dart' as app;
import '../test/test_environment.dart';

void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  late TestEnvironment env;

  setUpAll(() async {
    env = await TestEnvironment.create();
    await BJJModel.clearCachedLogin();
  });

  tearDownAll(() async {
    await env.tearDown();
  });

  Future<void> doSelectProjectNavBack(
    WidgetTester tester,
  ) async {
    final findBackButton = find.byKey(const Key('selectProjectNavBack'));
    await pumpUntilFound(tester, findBackButton);
    await tester.tap(findBackButton);
    await pumpUntilFound(tester, find.byTooltip('Sign up'));
  }

  Future<void> doSelectProject(
    WidgetTester tester,
    String name,
  ) async {
    final findProjectKey = find.byKey(Key('selectProject-$name'));
    await pumpUntilFound(tester, findProjectKey);
    await tester.tap(findProjectKey);
  }

  Future<void> doCreateProject(
    WidgetTester tester,
    TestProject project,
  ) async {
    final findNewProject = find.byKey(const Key('createNewProject'));
    await pumpUntilFound(tester, findNewProject);
    await tester.ensureVisible(findNewProject);
    await tester.pumpAndSettle();
    await tester.tap(findNewProject);
    await tester.pumpAndSettle();
    final findCreateProjectSubmit =
        find.byKey(const Key('createProjectSubmit'));
    await pumpUntilFound(tester, findCreateProjectSubmit);
    final findCreateProjectName = find.byKey(const Key('createProjectName'));
    await tester.enterText(findCreateProjectName, 'd');
    await tester.tap(findCreateProjectSubmit);
    await pumpUntilFound(tester,
        find.text('Must contain at least $minProjectNameLength characters'));
    await tester.enterText(findCreateProjectName, project.name);
    await tester.pumpAndSettle();
    final findSubmit = find.text('Create project "${project.name}"');
    await pumpUntilFound(tester, findSubmit);
    await tester.tap(findSubmit);
    await tester.pumpAndSettle();
    await pumpUntilFound(tester, find.byKey(const Key('accountTab')));
  }

  Future<void> doSignOut(WidgetTester tester) async {
    final findAccountTab = find.byKey(const Key('accountTab'));
    await pumpUntilFound(tester, findAccountTab);
    expect(findAccountTab, findsOneWidget);
    await tester.tap(findAccountTab);
    await tester.pumpAndSettle();
    await tester.tap(findAccountTab);
    await tester.pumpAndSettle();
    final findSignOut = find.byKey(const Key('signOut'));
    await pumpUntilFound(tester, findSignOut);
    await tester.tap(findSignOut);
    await pumpUntilFound(tester, find.byTooltip('Sign up'));
  }

  Future<void> maybeSignOut(WidgetTester tester) async {
    try {
      await doSignOut(tester);
    } catch (err) {
      // already signed out
    }
  }

  Future<void> doSignUp(WidgetTester tester, TestCredentials creds,
      {bool fumble = false}) async {
    await tester.pumpAndSettle();
    try {
      await pumpUntilFound(tester, find.byTooltip('Sign up'));
    } catch (err) {
      await maybeSignOut(tester);
    }
    await tester.tap(find.byTooltip('Sign up'));
    await tester.pumpAndSettle();
    final findSubmit = find.byTooltip('Create account');
    await pumpUntilFound(tester, findSubmit);
    expect(find.text('Create an Account'), findsOneWidget);
    final findUsername = find.byKey(const Key('signUpUsername'));
    expect(findUsername, findsOneWidget);
    final findEmail = find.byKey(const Key('signUpEmail'));
    expect(findEmail, findsOneWidget);
    final findFirstName = find.byKey(const Key('signUpFirstName'));
    expect(findFirstName, findsOneWidget);
    final findLastName = find.byKey(const Key('signUpLastName'));
    expect(findLastName, findsOneWidget);
    final findPassword = find.byKey(const Key('signUpPassword'));
    expect(findPassword, findsOneWidget);
    final findPasswordConfirm = find.byKey(const Key('signUpPasswordConfirm'));
    expect(findPasswordConfirm, findsOneWidget);
    if (fumble) {
      await tester.enterText(findUsername, 'a');
      await tester.pumpAndSettle();
      await tester.tap(findSubmit);
      await tester.pumpAndSettle();
      expect(find.text('Must contain at least $minPasswordLength characters'),
          findsOneWidget);
      await tester.enterText(findUsername, 'admin');
      await tester.tap(findSubmit);
      await pumpUntilFound(tester, find.text('Username already taken'));
      await tester.enterText(findUsername, creds.username);
      await tester.pumpAndSettle();
      await tester.enterText(findEmail, 'd@d.com');
      await pumpUntilFound(tester, find.text('Email is already taken'));
      await tester.enterText(findEmail, creds.email);
      await tester.enterText(findFirstName, creds.firstName);
      await tester.enterText(findLastName, creds.lastName);
      await tester.enterText(findPassword, 'D');
      await tester.pumpAndSettle();
      expect(find.text('Must contain at least $minPasswordLength characters'),
          findsOneWidget);
      await tester.enterText(findPassword, 'DDDDDDDD');
      await tester.pumpAndSettle();
      expect(find.text('Must contain lowercase letters'), findsOneWidget);
      await tester.enterText(findPassword, 'dddddddd');
      await tester.pumpAndSettle();
      expect(find.text('Must contain uppercase letters'), findsOneWidget);
      await tester.enterText(findPassword, 'ddddDDDD');
      await tester.pumpAndSettle();
      expect(find.text('Must contain numbers'), findsOneWidget);
      await tester.enterText(findPassword, 'ddddDDD1');
      await tester.pumpAndSettle();
      expect(find.text('Must contain special characters'), findsOneWidget);
      await tester.enterText(findPassword, creds.password);
      await tester.enterText(findPasswordConfirm, 'aaaaa');
      await tester.tap(findSubmit);
      await pumpUntilFound(tester, find.text('Passwords do not match'));
      expect(find.text('Passwords do not match'), findsOneWidget);
    } else {
      await tester.enterText(findUsername, creds.username);
      await tester.enterText(findEmail, creds.email);
      await tester.enterText(findFirstName, creds.firstName);
      await tester.enterText(findLastName, creds.lastName);
      await tester.enterText(findPassword, creds.password);
    }
    await tester.enterText(findPasswordConfirm, creds.password);
    await tester.pumpAndSettle();
    await tester.tap(findSubmit);
    await doSelectProjectNavBack(tester);
  }

  Future<void> doSignIn(WidgetTester tester, TestCredentials creds,
      {bool fumble = false}) async {
    await tester.pumpAndSettle();
    final findUsername = find.byKey(const Key('signInUsername'));
    await pumpUntilFound(tester, findUsername);
    expect(findUsername, findsOneWidget);
    final findPassword = find.byKey(const Key('signInPassword'));
    expect(findPassword, findsOneWidget);
    final findSubmit = find.byKey(const Key('signInSubmit'));
    expect(findSubmit, findsOneWidget);
    if (fumble) {
      await tester.tap(findSubmit);
      await pumpUntilFound(tester, find.text('Please enter your username'));
      expect(find.text('Please enter your password'), findsOneWidget);
      await tester.enterText(findUsername, creds.username);
      await tester.enterText(findPassword, 'fooBarBaz123%%');
      await tester.pumpAndSettle();
      await tester.tap(findSubmit);
      await pumpUntilFound(tester, find.text('Invalid username or password'));
    } else {
      await tester.enterText(findUsername, creds.username);
    }
    await tester.enterText(findPassword, creds.password);
    await tester.pumpAndSettle();
    await tester.tap(findSubmit);
    await tester.pumpAndSettle();
  }

  testWidgets('sign up', (WidgetTester tester) async {
    app.main();
    await doSignUp(tester, env.creds[0], fumble: true);
    await doSignUp(tester, env.creds[1]);
  });

  testWidgets('sign in', (WidgetTester tester) async {
    app.main();
    // the reason they're found is because they're in the widget tree
    await doSignIn(tester, env.creds[0], fumble: false);
    await doSelectProjectNavBack(tester);
  });

  testWidgets('create project', (WidgetTester tester) async {
    app.main();
    await doSignIn(tester, env.creds[0]);
    await doCreateProject(tester, env.projects[0]);
  });

  testWidgets('project name taken', (WidgetTester tester) async {
    app.main();
    await doSelectProject(tester, env.projects[0].name);
    // home screen is now visible
    final findAccountTab = find.byKey(const Key('accountTab'));
    await pumpUntilFound(tester, findAccountTab);
    await tester.tap(findAccountTab);
    await tester.pumpAndSettle();
    final findSwitchProject = find.byKey(const Key('switchProject'));
    await pumpUntilFound(tester, findSwitchProject);
    await tester.tap(findSwitchProject);
    final findNewProject = find.byKey(const Key('createNewProject'));
    await pumpUntilFound(tester, findNewProject);

    // switch projects page is now visible
    await tester.tap(findNewProject);
    final findCreateProjectName = find.byKey(const Key('createProjectName'));
    await pumpUntilFound(tester, findCreateProjectName);
    await tester.enterText(findCreateProjectName, env.projects[0].name);
    await pumpUntilFound(
        tester, find.text('A project with this name already exists'));
  });

  testWidgets('add video to project', (WidgetTester tester) async {
    app.main();
    final findVideos = find.byKey(const Key('videos'));
    await pumpUntilFound(tester, findVideos);
    await tester.tap(findVideos);
    final findAddVideo = find.byKey(const Key('addVideo'));
    await pumpUntilFound(tester, findAddVideo);
    await tester.tap(findAddVideo);
    final findAddVideoInput = find.byKey(const Key('addVideoInput'));
    await pumpUntilFound(tester, findAddVideoInput);
    await tester.enterText(findAddVideoInput, env.projects[0].inputVideos[0]);
    final findConfirmAddVideo = find.byKey(const Key('confirmAddVideo'));
    await tester.tap(findConfirmAddVideo);
    final findVideo =
        find.byKey(Key('video-${env.projects[0].inputVideos[0]}'));
    await pumpUntilFound(tester, findVideo);
    await tester.tap(find.byKey(const Key('videosNavBack')));
    final findOutputTab = find.byKey(const Key('outputTab'));
    await pumpUntilFound(tester, findOutputTab);
    await tester.tap(findOutputTab);
    await pumpUntilFound(tester, findVideo);
  });

  testWidgets('add playlist to project', (WidgetTester tester) async {
    app.main();
    final findPlaylists = find.byKey(const Key('playlists'));
    await pumpUntilFound(tester, findPlaylists);
    await tester.tap(findPlaylists);
    final findAddPlaylist = find.byKey(const Key('addPlaylist'));
    await pumpUntilFound(tester, findAddPlaylist);
    await tester.tap(findAddPlaylist);
    final findAddPlaylistInput = find.byKey(const Key('addPlaylistInput'));
    await pumpUntilFound(tester, findAddPlaylistInput);
    await tester.enterText(
        findAddPlaylistInput, env.projects[0].inputPlaylists[0]);
    final findConfirmAddPlaylist = find.byKey(const Key('confirmAddPlaylist'));
    await tester.tap(findConfirmAddPlaylist);
    final findPlaylist =
        find.byKey(Key('playlist-${env.projects[0].inputPlaylists[0]}'));
    await pumpUntilFound(tester, findPlaylist);
    await tester.tap(find.byKey(const Key('playlistsNavBack')));
  });

  testWidgets('add channel to project', (WidgetTester tester) async {
    app.main();
    final findChannels = find.byKey(const Key('channels'));
    await pumpUntilFound(tester, findChannels);
    await tester.tap(findChannels);
    final findAddChannel = find.byKey(const Key('addChannel'));
    await pumpUntilFound(tester, findAddChannel);
    await tester.tap(findAddChannel);
    final findAddChannelInput = find.byKey(const Key('addChannelInput'));
    await pumpUntilFound(tester, findAddChannelInput);
    await tester.enterText(
        findAddChannelInput, env.projects[0].inputChannels[0]);
    final findConfirmAddChannel = find.byKey(const Key('confirmAddChannel'));
    await tester.tap(findConfirmAddChannel);
    final findChannel =
        find.byKey(Key('channel-${env.projects[0].inputChannels[0]}'));
    await pumpUntilFound(tester, findChannel);
    await tester.tap(find.byKey(const Key('channelsNavBack')));
    await doSignOut(tester);
  });

  testWidgets('ensure project is private', (WidgetTester tester) async {
    app.main();
    await doSignIn(tester, env.creds[1]);
    await pumpUntilFound(tester, find.byKey(const Key('createNewProject')));
    final findProjectKey =
        find.byKey(Key('selectProject-${env.projects[0].name}'));
    expect(findProjectKey, findsNothing);
    await doSelectProjectNavBack(tester);
  });

  testWidgets('add collaborator', (WidgetTester tester) async {
    app.main();
    await doSignIn(tester, env.creds[0]);
    await doSelectProject(tester, env.projects[0].name);
  });
}

Future<void> pumpUntilFound(
  WidgetTester tester,
  Finder finder, {
  Duration timeout = const Duration(seconds: 8),
}) async {
  bool timerDone = false;
  final timer = Timer(timeout, () => timerDone = true);
  while (timerDone != true) {
    await tester.pump();

    final found = tester.any(finder);
    if (found) {
      timerDone = true;
    }
  }
  timer.cancel();
  expect(finder, findsWidgets);
}
