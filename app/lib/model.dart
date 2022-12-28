import 'package:flutter/material.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:scoped_model/scoped_model.dart';
import 'api.dart' as api;

// https://img.youtube.com/vi/e5YuPpbzBdo/maxresdefault.jpg

const credStorageProjectId = 't4vdp';
const credStorageUsername = 't4vdcu';
const credStoragePassword = 't4vdcp';
const credStorageBrightness = 't4vddm';

final List<api.VideoListItem> exampleVideos = [
  api.VideoListItem(
    id: "6wrYM8KzBRU",
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
  api.VideoListItem(
    id: "e5YuPpbzBdo",
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
];

final List<api.PlaylistListItem> examplePlaylists = [
  api.PlaylistListItem(
    id: "PLuWwmKO5nQLtJ7aczfnESU3Wvetw-N12I",
    channel: "ROYDEAN",
    channelId: "ROYDEAN",
    title: "Costa Rica 2021",
    numVideos: 6,
  ),
];

final List<api.ChannelListItem> exampleChannels = [
  api.ChannelListItem(
    id: "tylerspangler",
    name: "Tyler Spangler",
    avatarUrl: "assets/tylerspangler.jpg",
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
  List<api.PlaylistListItem> get playlists => _playlists;
  List<api.ChannelListItem> get channels => _channels;
  List<api.VideoListItem> get videos => _videos;
  api.Dataset? get dataset => _dataset;
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
  api.Project? get project => _project;

  List<api.Project> _projects = [];
  List<api.PlaylistListItem> _playlists = [];
  List<api.ChannelListItem> _channels = [];
  List<api.VideoListItem> _videos = [];
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
    BuildContext context,
    api.SearchUser user,
  ) async {
    await withAuth(context, () async {
      await ensureProject(context);
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
    BuildContext context,
    String userId,
  ) async {
    await withAuth(context, () async {
      await ensureProject(context);
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
          BuildContext context, String prefix) async =>
      await withAuth(context,
          () async => await api.searchUsers(creds: _creds!, prefix: prefix));

  Future<void> createProject({
    required BuildContext context,
    required String name,
  }) async {
    await withAuth(context, () async {
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
    }).catchError((err) {
      _loginErr = err is api.InvalidCredentialsError
          ? 'Invalid username or password'
          : err.toString();
      notifyListeners();
    });
  }

  Future<void> signOut() async {
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

  Future<bool> ensureCreds(BuildContext context, [mounted = true]) async {
    if (isLoggedIn) return true;
    await readCachedCreds();
    if (isLoggedIn) return true;
    if (!mounted) return false;
    Navigator.of(context).pushNamed('/splash');
    return false;
  }

  Future<dynamic> withAuth(BuildContext context, Future<dynamic> Function() f,
      [mounted = true]) async {
    if (!await ensureCreds(context)) {
      await Navigator.of(context).pushNamed('/splash');
      if (!isLoggedIn) throw UnimplementedError('re-login failed');
    }
    try {
      return await f();
    } on api.InvalidCredentialsError {
      if (!mounted) return;
      await Navigator.of(context).pushNamed('/splash');
      return await f(); // try a second time without catching
    }
  }

  Future<void> refreshProjects(BuildContext context) async {
    await withAuth(context, () async {
      _projects = await api.listProjects(creds: _creds!);
      notifyListeners();
    });
  }

  Future<void> selectProject(
    BuildContext context,
    String projectId,
  ) async {
    await withAuth(context, () async {
      _project = await api.getProject(creds: _creds!, projectId: projectId);
      await const FlutterSecureStorage().write(
        key: credStorageProjectId,
        value: _project!.id,
      );
      notifyListeners();
    });
  }

  Future<void> refreshDataset(BuildContext context) async {
    await withAuth(context, () async {
      await ensureProject(context);
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
      } on api.ResourceNotFoundError {
        // project was deleted
        await const FlutterSecureStorage().delete(key: credStorageProjectId);
        return;
      }
      notifyListeners();
    }
  }

  Future<void> ensureProject(
    BuildContext context, {
    mounted = true,
  }) async {
    if (hasProject) return;
    await loadCachedProject();
    if (hasProject) return;
    if (!mounted) return;
    await Navigator.of(context).pushNamed('/splash');
  }

  Future<void> refreshMarkers(BuildContext context) async {
    await withAuth(context, () async {
      await ensureProject(context);
      _markers = await api.getStack(
        projectId: _project!.id,
        creds: creds!,
      );
      _markerIndex = 0;
      notifyListeners();
    });
  }

  Future<void> classify({
    required BuildContext context,
    required bool label,
  }) async {
    await withAuth(context, () async {
      await ensureProject(context);
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
    required BuildContext context,
    required String input,
    required bool blacklist,
  }) async {
    await withAuth(context, () async {
      await ensureProject(context);
      var item = await api.addChannel(
        projectId: _project!.id,
        creds: _creds!,
        input: input,
        blacklist: blacklist,
      );
      _channels.add(item);
      notifyListeners();
    });
  }

  Future<void> addPlaylist({
    required BuildContext context,
    required String input,
    required bool blacklist,
  }) async {
    await withAuth(context, () async {
      await ensureProject(context);
      var item = await api.addPlaylist(
        projectId: _project!.id,
        creds: _creds!,
        input: input,
        blacklist: blacklist,
      );
      _playlists.add(item);
      notifyListeners();
    });
  }

  Future<void> addVideo({
    required BuildContext context,
    required String input,
    required bool blacklist,
  }) async {
    await withAuth(context, () async {
      await ensureProject(context);
      var item = await api.addVideo(
        projectId: _project!.id,
        creds: _creds!,
        input: input,
        blacklist: blacklist,
      );
      _videos.add(item);
      notifyListeners();
    });
  }

  Future<void> refreshChannels(BuildContext context) async {
    await withAuth(context, () async {
      await ensureProject(context);
      _channels = await api.listChannels(
        projectId: _project!.id,
        creds: _creds!,
      );
      notifyListeners();
    });
  }

  Future<void> refreshPlaylists(BuildContext context) async {
    await withAuth(context, () async {
      await ensureProject(context);
      _playlists = await api.listPlaylists(
        projectId: _project!.id,
        creds: _creds!,
      );
      notifyListeners();
    });
  }

  Future<void> refreshVideos(BuildContext context) async {
    await withAuth(context, () async {
      await ensureProject(context);
      _videos = await api.listVideos(
        projectId: _project!.id,
        creds: _creds!,
      );
      notifyListeners();
    });
  }

  Future<void> removeChannel({
    required BuildContext context,
    required String id,
    required bool blacklist,
  }) async {
    await withAuth(context, () async {
      await ensureProject(context);
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
    required BuildContext context,
    required String id,
    required bool blacklist,
  }) async {
    await withAuth(context, () async {
      await ensureProject(context);
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
    required BuildContext context,
    required String id,
    required bool blacklist,
  }) async {
    await withAuth(context, () async {
      await ensureProject(context);
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
