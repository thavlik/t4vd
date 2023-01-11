import 'dart:convert';

import 'package:web_socket_channel/io.dart';
import 'package:flutter/material.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:scoped_model/scoped_model.dart';
import 'api.dart' as api;

// https://img.youtube.com/vi/e5YuPpbzBdo/maxresdefault.jpg

const credStorageProjectId = 't4vdp';
const credStorageUsername = 't4vdcu';
const credStoragePassword = 't4vdcp';
const credStorageBrightness = 't4vddm';

final List<api.Video> exampleVideos = [
  api.Video(
    id: "6wrYM8KzBRU",
    info: api.VideoInfo(
      channel: "Tyler Spangler",
      channelId: "TylerSpangler",
      title: "I Tried Using Jiu Jitsu Versus A Knife",
      uploadDate: "20221128",
      fps: 60,
      duration: 495,
      width: 1920,
      height: 1080,
      thumbnail: "assets/tylerspangler.jpg",
      uploader: "Tyler Spangler",
      uploaderId: "@tylerspanger",
    ),
  ),
  api.Video(
    id: "e5YuPpbzBdo",
    info: api.VideoInfo(
      channel: "Tyler Spangler",
      channelId: "TylerSpangler",
      title: "I Tore My Knee At A New Gym",
      uploadDate: "20221122",
      fps: 60,
      duration: 495,
      width: 1920,
      height: 1080,
      thumbnail: "assets/tylerspangler.jpg",
      uploader: "Tyler Spangler",
      uploaderId: "@tylerspanger",
    ),
  ),
];

final List<api.Playlist> examplePlaylists = [
  api.Playlist(
    id: "PLuWwmKO5nQLtJ7aczfnESU3Wvetw-N12I",
    info: api.PlaylistInfo(
      channel: "ROYDEAN",
      channelId: "ROYDEAN",
      title: "Costa Rica 2021",
      numVideos: 6,
    ),
  ),
];

final List<api.Channel> exampleChannels = [
  api.Channel(
    id: "tylerspangler",
    info: api.ChannelInfo(
      name: "Tyler Spangler",
      avatarUrl: "assets/tylerspangler.jpg",
    ),
  )
];

final List<api.Project> exampleProjects = [
  api.Project(
    id: 'example-0',
    name: 'Brazilian Jiu-jitsu Vision',
  ),
  api.Project(
    id: 'example-1',
    name: 'Doom Vision',
  ),
  api.Project(
    id: 'example-2',
    name: 'Quake Vision',
  )
];

class UserCredentials {
  final String id;
  final String username;
  final String email;
  final String firstName;
  final String lastName;
  final String accessToken;
  final bool enabled;

  UserCredentials({
    required this.id,
    required this.username,
    required this.email,
    required this.firstName,
    required this.lastName,
    required this.accessToken,
    required this.enabled,
  });

  static UserCredentials fromMap(Map<dynamic, dynamic> m) => UserCredentials(
        id: m['id'],
        username: m['username'],
        email: m['email'],
        firstName: m['firstName'],
        lastName: m['lastName'],
        accessToken: m['accessToken'],
        enabled: m['enabled'],
      );

  Map<String, dynamic> toMap() => {
        'id': id,
        'username': username,
        'email': email,
        'firstName': firstName,
        'lastName': lastName,
        'accessToken': accessToken,
        'enabled': enabled,
      };
}

class BJJModel extends Model {
  BJJModel({UserCredentials? creds}) : _creds = creds;

  Brightness _brightness = Brightness.dark;

  Brightness get brightness => _brightness;

  Future<void> setBrightness(Brightness value) async {
    _brightness = value;
    const FlutterSecureStorage().write(
        key: credStorageBrightness,
        value: value == Brightness.dark ? 'dark' : 'light');
    notifyListeners();
  }

  List<api.Project> get projects => _projects;
  List<api.Playlist> get playlists => _playlists;
  List<api.Channel> get channels => _channels;
  List<api.Video> get videos => _videos;
  List<api.Marker> get markers => _markers;
  int get markerIndex => _markerIndex;
  String? get loginErr => _loginErr;
  UserCredentials? get creds => _creds;
  api.Marker? get currentMarker =>
      _markers.isNotEmpty && _markerIndex < _markers.length
          ? _markers[_markerIndex]
          : null;
  bool get isLoggedIn => _creds != null;
  bool get hasProject => _project != null;
  String? get projectId => _project?.id;
  api.Dataset? get dataset => _dataset;
  api.Project? get project => _project;

  List<api.Project> _projects = [];
  List<api.Playlist> _playlists = [];
  List<api.Channel> _channels = [];
  List<api.Video> _videos = [];
  List<api.Marker> _markers = [];
  int _markerIndex = 0;
  api.Dataset? _dataset;
  String? _loginErr;
  UserCredentials? _creds;
  api.Project? _project;

  Future<void> register({
    required String username,
    required String email,
    required String firstName,
    required String lastName,
    required String password,
  }) async {
    _creds = await api.createAccount(
      username: username,
      email: email,
      firstName: firstName,
      lastName: lastName,
      password: password,
    );
    notifyListeners();
  }

  Future<void> addCollaborator(
    NavigatorState nav,
    api.SearchUser user,
  ) async {
    await withAuth(nav, () async {
      await ensureProject(nav);
      await api.addCollaborator(
        creds: _creds!,
        userId: user.id,
        projectId: _project!.id,
      );
      _project!.collaborators ??= [];
      if (_project!.collaborators!.indexWhere((u) => u.id == user.id) == -1) {
        _project!.collaborators!.add(user);
      }
      notifyListeners();
    });
  }

  Future<void> removeCollaborator(
    NavigatorState nav,
    String userId,
  ) async {
    await withAuth(nav, () async {
      await ensureProject(nav);
      await api.removeCollaborator(
        creds: _creds!,
        userId: userId,
        projectId: _project!.id,
      );
      _project!.collaborators?.removeWhere((user) => user.id == userId);
      notifyListeners();
    });
  }

  Future<List<api.SearchUser>> searchUsers(
          NavigatorState nav, String prefix) async =>
      await withAuth(nav,
          () async => await api.searchUsers(creds: _creds!, prefix: prefix));

  Future<void> createProject({
    required NavigatorState nav,
    required String name,
  }) async {
    await withAuth(nav, () async {
      final project = await api.createProject(
        creds: _creds!,
        name: name,
      );
      _projects.add(project);
      _project = project;
      notifyListeners();
    });
  }

  static Future<void> clearCachedLogin() async => await Future.wait([
        const FlutterSecureStorage().delete(key: credStorageUsername),
        const FlutterSecureStorage().delete(key: credStoragePassword),
        const FlutterSecureStorage().delete(key: credStorageProjectId),
        const FlutterSecureStorage().delete(key: credStorageBrightness),
      ]);

  Future<void> readCachedCreds() async {
    final username =
        await const FlutterSecureStorage().read(key: credStorageUsername);
    final password =
        await const FlutterSecureStorage().read(key: credStoragePassword);
    if (username == null || password == null) return;
    await login(username: username, password: password);
  }

  Future<void> readCachedBrightness() async {
    final value =
        await const FlutterSecureStorage().read(key: credStorageBrightness);
    if (value == 'dark' || value == null) {
      _brightness = Brightness.dark;
    } else if (value == 'light') {
      _brightness = Brightness.light;
    } else {
      throw ErrorSummary('invalid cached brightness value "$value"');
    }
    notifyListeners();
  }

  Future<void> writeCachedCreds({
    required String username,
    required String password,
  }) async {
    await const FlutterSecureStorage()
        .write(key: credStorageUsername, value: username);
    await const FlutterSecureStorage()
        .write(key: credStoragePassword, value: password);
  }

  Future<void> login({
    required String username,
    required String password,
  }) async {
    _loginErr = null;
    notifyListeners();
    await api.login(username, password).then((value) async {
      _creds = value;
      _loginErr = null;
      notifyListeners();
      await writeCachedCreds(
        username: username,
        password: password,
      );
      await connectWebSock();
    }).catchError((err) {
      _loginErr = err is api.InvalidCredentialsError
          ? 'Invalid username or password'
          : err.toString();
      notifyListeners();
    });
  }

  Future<void> signOut() async {
    _channel?.sink.close();
    final c = clearCachedLogin();
    _project = null;
    if (!isLoggedIn) {
      await c;
      return;
    }
    final f = Future.wait([
      c,
      api.signOut(creds!),
    ]);
    _creds = null;
    await f;
    notifyListeners();
  }

  IOWebSocketChannel? _channel;

  Future<void> connectWebSock() async {
    if (!isLoggedIn) throw ErrorSummary('not logged in');
    _channel?.sink.close();
    _channel = api.connectWebSock(_creds!);
    _channel!.stream.handleError((obj, stackTrace) {
      // reconnect if we're still logged in
      print('websock error: $obj');
      print(stackTrace);
      if (_creds == null) return;
      connectWebSock();
    });
    _channel!.stream.listen((message) => handleWebSockMessage(message));
    if (_project != null) {
      _channel!.sink.add(jsonEncode({
        'type': 'subscribe',
        'projectID': _project!.id,
        'unsubscribeAll': true,
      }));
    }
  }

  void handleWebSockMessage(dynamic message) {
    print(message);
    final obj = jsonDecode(message) as Map<String, dynamic>;
    switch (obj['type']) {
      case 'channel_details':
      case 'playlists_details':
      case 'video_details':
      case 'channel_video':
      case 'playlist_video':
    }
  }

  void precacheFrames(BuildContext context) {
    try {
      for (var marker in _markers.sublist(_markerIndex)) {
        precacheImage(
          NetworkImage(api.videoThumbnail(marker.videoId)),
          context,
        );
      }
    } catch (err) {
      // do nothing
    }
  }

  Future<bool> ensureCreds(NavigatorState nav) async {
    if (isLoggedIn) return true;
    await readCachedCreds();
    if (isLoggedIn) return true;
    nav.pushNamed('/splash');
    return false;
  }

  Future<dynamic> withAuth(
    NavigatorState nav,
    Future<dynamic> Function() f,
  ) async {
    if (!await ensureCreds(nav)) {
      await nav.pushNamed('/splash');
      if (!isLoggedIn) throw UnimplementedError('re-login failed');
    }
    try {
      return await f();
    } on api.InvalidCredentialsError {
      await nav.pushNamed('/splash');
      return await f(); // try a second time without catching
    } on api.ForbiddenError {
      await nav.pushNamed('/splash');
      return await f(); // try a second time without catching
    }
  }

  Future<void> refreshProjects(NavigatorState nav) async {
    await withAuth(nav, () async {
      _projects = await api.listProjects(creds: _creds!);
      notifyListeners();
    });
  }

  Future<void> selectProject(
    NavigatorState nav,
    String projectId,
  ) async {
    await withAuth(nav, () async {
      _project = await api.getProject(creds: _creds!, projectId: projectId);
      await const FlutterSecureStorage().write(
        key: credStorageProjectId,
        value: _project!.id,
      );
      if (_channel != null) {
        _channel!.sink.add(jsonEncode({
          'type': 'subscribe',
          'projectID': _project!.id,
          'unsubscribeAll': true,
        }));
      }
      notifyListeners();
    });
  }

  Future<void> refreshDataset(NavigatorState nav) async {
    await withAuth(nav, () async {
      await ensureProject(nav);
      _dataset = await api.getDataset(
        projectId: _project!.id,
        creds: _creds!,
      );
      notifyListeners();
    });
  }

  Future<void> loadCachedProject() async {
    final projectId =
        await const FlutterSecureStorage().read(key: credStorageProjectId);
    if (projectId != null) {
      try {
        _project = await api.getProject(
          projectId: projectId,
          creds: _creds!,
        );
      } on api.ForbiddenError {
        // user no longer has access to project
        await const FlutterSecureStorage().delete(key: credStorageProjectId);
        return;
      } on api.ResourceNotFoundError {
        // project was deleted
        await const FlutterSecureStorage().delete(key: credStorageProjectId);
        return;
      }
      notifyListeners();
    }
  }

  Future<void> ensureProject(NavigatorState nav) async {
    if (hasProject) return;
    await loadCachedProject();
    if (hasProject) return;
    await nav.pushNamed('/splash');
  }

  Future<void> refreshMarkers(NavigatorState nav) async {
    await withAuth(nav, () async {
      await ensureProject(nav);
      _markers = [
        await api.getRandomMarker(
          projectId: _project!.id,
          creds: creds!,
        )
      ];
      _markerIndex = 0;
      notifyListeners();
    });
  }

  Future<void> classify({
    required NavigatorState nav,
    required bool label,
  }) async {
    await withAuth(nav, () async {
      await ensureProject(nav);
      final cur = _markers[_markerIndex];
      await api.classifyMarker(
        projectId: _project!.id,
        videoId: cur.videoId,
        time: cur.time,
        label: label,
        creds: creds!,
      );
      _markerIndex++;
      final remaining = _markers.length - _markerIndex;
      if (remaining < 3) {
        // asynchronously get another stack
        api.getStack(projectId: _project!.id, creds: creds!).then((value) {
          _markers.addAll(value);
          notifyListeners();
        });
      }
      notifyListeners();
    });
  }

  Future<void> classifyBack() async {
    _markerIndex--;
    notifyListeners();
  }

  Future<void> addChannel({
    required NavigatorState nav,
    required String input,
    required bool blacklist,
  }) async {
    await withAuth(nav, () async {
      await ensureProject(nav);
      var item = await api.addChannel(
        projectId: _project!.id,
        creds: _creds!,
        input: input,
        blacklist: blacklist,
      );
      if (item != null) {
        _channels.add(item);
        notifyListeners();
      }
    });
  }

  Future<void> addPlaylist({
    required NavigatorState nav,
    required String input,
    required bool blacklist,
  }) async {
    await withAuth(nav, () async {
      await ensureProject(nav);
      var item = await api.addPlaylist(
        projectId: _project!.id,
        creds: _creds!,
        input: input,
        blacklist: blacklist,
      );
      if (item != null) {
        _playlists.add(item);
        notifyListeners();
      }
    });
  }

  Future<void> addVideo({
    required NavigatorState nav,
    required String input,
    required bool blacklist,
  }) async {
    await withAuth(nav, () async {
      await ensureProject(nav);
      var item = await api.addVideo(
        projectId: _project!.id,
        creds: _creds!,
        input: input,
        blacklist: blacklist,
      );
      if (item != null) {
        _videos.add(item);
        notifyListeners();
      }
    });
  }

  Future<void> refreshChannels(NavigatorState nav) async {
    await withAuth(nav, () async {
      await ensureProject(nav);
      _channels = await api.listChannels(
        projectId: _project!.id,
        creds: _creds!,
      );
      notifyListeners();
    });
  }

  Future<void> refreshPlaylists(NavigatorState nav) async {
    await withAuth(nav, () async {
      await ensureProject(nav);
      _playlists = await api.listPlaylists(
        projectId: _project!.id,
        creds: _creds!,
      );
      notifyListeners();
    });
  }

  Future<void> refreshVideos(NavigatorState nav) async {
    await withAuth(nav, () async {
      await ensureProject(nav);
      _videos = await api.listVideos(
        projectId: _project!.id,
        creds: _creds!,
      );
      notifyListeners();
    });
  }

  Future<void> removeChannel({
    required NavigatorState nav,
    required String id,
    required bool blacklist,
  }) async {
    await withAuth(nav, () async {
      await ensureProject(nav);
      await api.removeChannel(
        projectId: _project!.id,
        creds: _creds!,
        id: id,
        blacklist: blacklist,
      );
      _channels.removeWhere((a) => a.id == id);
      notifyListeners();
    });
  }

  Future<void> removePlaylist({
    required NavigatorState nav,
    required String id,
    required bool blacklist,
  }) async {
    await withAuth(nav, () async {
      await ensureProject(nav);
      await api.removePlaylist(
        projectId: _project!.id,
        creds: _creds!,
        id: id,
        blacklist: blacklist,
      );
      _playlists.removeWhere((a) => a.id == id);
      notifyListeners();
    });
  }

  Future<void> removeVideo({
    required NavigatorState nav,
    required String id,
    required bool blacklist,
  }) async {
    await withAuth(nav, () async {
      await ensureProject(nav);
      await api.removeVideo(
        projectId: _project!.id,
        id: id,
        blacklist: blacklist,
        creds: creds!,
      );
      _videos.removeWhere((a) => a.id == id);
      notifyListeners();
    });
  }
}
