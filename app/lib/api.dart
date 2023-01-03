import 'dart:convert';

import 'package:t4vd/model.dart';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import 'package:http/http.dart';

const apiHost = 'api.beebs.dev';
const int minUsernameLength = 4;
const int minProjectNameLength = 3;
const int minPasswordLength = 8;

class InvalidCredentialsError extends Error {}

class ResourceNotFoundError extends Error {}

class Project {
  final String id;
  final String name;
  List<SearchUser>? collaborators;

  Project({
    required this.id,
    required this.name,
    this.collaborators,
  });

  static Project fromMap(Map<dynamic, dynamic> m) => Project(
        id: m['id'],
        name: m['name'],
        collaborators: (m['collaborators'] as List?)
            ?.map((e) => SearchUser.fromMap(e))
            .toList(),
      );
}

class SearchUser {
  final String id;
  final String username;

  SearchUser({
    required this.id,
    required this.username,
  });

  static SearchUser fromMap(Map<dynamic, dynamic> m) => SearchUser(
        id: m['id'],
        username: m['username'],
      );
}

class Marker {
  final String videoId;
  final int time;

  Marker({
    required this.videoId,
    required this.time,
  });

  String get imageUrl => "https://$apiHost/frame?v=$videoId&t=$time";

  static Marker fromMap(Map m) => Marker(
        videoId: m['videoID'],
        time: m['time'],
      );
}

class Dataset {
  final String id;
  final List<Video> videos;
  final DateTime timestamp;

  Dataset({
    required this.id,
    required this.videos,
    required this.timestamp,
  });

  static Dataset fromMap(Map m) => Dataset(
        id: m['id'],
        timestamp: DateTime.fromMicrosecondsSinceEpoch(m['timestamp'] ~/ 1000),
        videos: (m['videos'] as List).map((o) => Video.fromMap(o)).toList(),
      );
}

enum VideoProgressState {
  pendingQuery,
  querying,
  pendingDownload,
  downloading,
}

class VideoProgress {
  VideoProgressState state;

  VideoProgress({
    required this.state,
  });

  static VideoProgress fromMap(Map<dynamic, dynamic> m) => VideoProgress(
      state:
          VideoProgressState.values.firstWhere((state) => state == m['state']));
}

class VideoInfo {
  final String title;
  final String uploader;
  final String uploaderId;
  final String channel;
  final String channelId;
  final String thumbnail;
  final int duration;
  final int width;
  final int height;
  final String uploadDate;
  final int fps;

  VideoInfo({
    required this.uploader,
    required this.uploaderId,
    required this.thumbnail,
    required this.title,
    required this.channel,
    required this.channelId,
    required this.duration,
    required this.width,
    required this.height,
    required this.uploadDate,
    required this.fps,
  });

  static VideoInfo fromMap(Map m) => VideoInfo(
        title: m['title'],
        channel: m['channel'],
        channelId: m['channelID'],
        uploader: m['uploader'],
        uploaderId: m['uploaderID'],
        thumbnail: m['thumbnail'],
        duration: m['duration'],
        width: m['width'],
        height: m['height'],
        uploadDate: m['uploadDate'],
        fps: m['fPS'],
      );
}

class Video {
  final String id;
  bool blacklist;
  VideoInfo? info;
  VideoProgress? progress;

  Video({
    required this.id,
    this.blacklist = false,
    this.info,
    this.progress,
  });

  static Video fromMap(Map m) => Video(
        id: m['id'],
        blacklist: m['blacklist'] ?? false,
        info: m.containsKey('info') ? VideoInfo.fromMap(m['info']) : null,
        progress: m.containsKey('progress')
            ? VideoProgress.fromMap(m['progress'])
            : null,
      );
}

class PlaylistInfo {
  final String title;
  final String channel;
  final String channelId;
  final int numVideos;

  PlaylistInfo({
    required this.title,
    required this.channel,
    required this.channelId,
    required this.numVideos,
  });

  static PlaylistInfo fromMap(Map m) => PlaylistInfo(
        title: m['title'],
        channel: m['channel'],
        channelId: m['channelID'],
        numVideos: m['numVideos'],
      );
}

class PlaylistProgress {
  static PlaylistProgress fromMap(Map m) => PlaylistProgress();
}

class Playlist {
  final String id;
  bool blacklist;
  PlaylistInfo? info;
  PlaylistProgress? progress;

  Playlist({
    required this.id,
    this.blacklist = false,
    this.info,
    this.progress,
  });

  static Playlist fromMap(Map m) => Playlist(
        id: m['id'],
        blacklist: m['blacklist'],
        info: m.containsKey('info') ? PlaylistInfo.fromMap(m['info']) : null,
        progress: m.containsKey('progress')
            ? PlaylistProgress.fromMap(m['progress'])
            : null,
      );
}

class ChannelProgress {
  static ChannelProgress fromMap(Map m) => ChannelProgress();
}

class ChannelInfo {
  final String name;
  final String avatarUrl;

  ChannelInfo({
    required this.name,
    required this.avatarUrl,
  });

  static ChannelInfo fromMap(Map m) => ChannelInfo(
        name: m['name'],
        avatarUrl: m['avatar'],
      );
}

class Channel {
  final String id;
  bool blacklist;
  ChannelInfo? info;
  ChannelProgress? progress;

  Channel({
    required this.id,
    this.blacklist = false,
    this.info,
    this.progress,
  });

  static Channel fromMap(Map m) => Channel(
        id: m['id'],
        blacklist: m['blacklist'],
        info: m.containsKey('info') ? ChannelInfo.fromMap(m['info']) : null,
        progress: m.containsKey('progress')
            ? ChannelProgress.fromMap(m['progress'])
            : null,
      );
}

class SourceInput {
  List<Video> videos = [];
  List<Playlist> playlists = [];
  List<Channel> channels = [];
}

class SourceOutput {
  List<Video> videos = [];
}

// https://img.youtube.com/vi/e5YuPpbzBdo/maxresdefault.jpg

void checkHttpStatus(Response response) {
  if (response.statusCode == 401) {
    throw InvalidCredentialsError();
  } else if (response.statusCode == 404) {
    throw ResourceNotFoundError();
  } else if (response.statusCode != 200 && response.statusCode != 202) {
    throw ErrorSummary("status ${response.statusCode}: ${response.body}");
  }
}

Future<Dataset> getDataset({
  required String projectId,
  required UserCredentials creds,
}) async {
  final url = Uri.https(apiHost, 'dataset', {
    'p': projectId,
  });
  final response = await http.get(url, headers: {
    'AccessToken': creds.accessToken,
  });
  checkHttpStatus(response);
  final decodedResponse = jsonDecode(utf8.decode(response.bodyBytes)) as Map;
  return Dataset.fromMap(decodedResponse);
}

Future<Project> getProject({
  required String projectId,
  required UserCredentials creds,
}) async {
  final url = Uri.https(apiHost, 'project', {
    'p': projectId,
  });
  final response = await http.get(url, headers: {
    'AccessToken': creds.accessToken,
  });
  checkHttpStatus(response);
  final decodedResponse = jsonDecode(utf8.decode(response.bodyBytes)) as Map;
  return Project.fromMap(decodedResponse);
}

Future<Channel?> addChannel({
  required String projectId,
  required String input,
  required bool blacklist,
  required UserCredentials creds,
}) async {
  final url = Uri.https(apiHost, 'channel/add');
  final response = await http.post(url,
      headers: {'AccessToken': creds.accessToken},
      body: json.encode({
        'projectID': projectId,
        'input': input,
        'blacklist': blacklist,
      }));
  checkHttpStatus(response);
  if (response.statusCode == 202) {
    // channel info was not cached but the query was scheduled
    return null;
  }
  final decodedResponse = jsonDecode(utf8.decode(response.bodyBytes)) as Map;
  return Channel.fromMap(decodedResponse);
}

Future<Playlist?> addPlaylist({
  required String projectId,
  required String input,
  required bool blacklist,
  required UserCredentials creds,
}) async {
  final url = Uri.https(apiHost, 'playlist/add');
  final response = await http.post(url,
      headers: {'AccessToken': creds.accessToken},
      body: json.encode({
        'projectID': projectId,
        'input': input,
        'blacklist': blacklist,
      }));
  checkHttpStatus(response);
  if (response.statusCode == 202) {
    // channel info was not cached but the query was scheduled
    return null;
  }
  final decodedResponse = jsonDecode(utf8.decode(response.bodyBytes)) as Map;
  return Playlist.fromMap(decodedResponse);
}

Future<Video?> addVideo({
  required String projectId,
  required String input,
  required bool blacklist,
  required UserCredentials creds,
}) async {
  final url = Uri.https(apiHost, 'video/add');
  final response = await http.post(url,
      headers: {'AccessToken': creds.accessToken},
      body: json.encode({
        'projectID': projectId,
        'input': input,
        'blacklist': blacklist,
      }));
  checkHttpStatus(response);
  if (response.statusCode == 202) {
    // channel info was not cached but the query was scheduled
    return null;
  }
  final decodedResponse = jsonDecode(utf8.decode(response.bodyBytes)) as Map;
  return Video.fromMap(decodedResponse);
}

Future<List<Channel>> listChannels({
  required String projectId,
  required UserCredentials creds,
}) async {
  final url = Uri.https(apiHost, 'channel/list', {
    'p': projectId,
  });
  final response = await http.get(url, headers: {
    'AccessToken': creds.accessToken,
  });
  checkHttpStatus(response);
  final decodedResponse =
      jsonDecode(utf8.decode(response.bodyBytes)) as List? ?? [];
  return decodedResponse.map((e) => Channel.fromMap(e)).toList();
}

Future<List<Playlist>> listPlaylists({
  required String projectId,
  required UserCredentials creds,
}) async {
  final url = Uri.https(apiHost, 'playlist/list', {
    'p': projectId,
  });
  final response = await http.get(url, headers: {
    'AccessToken': creds.accessToken,
  });
  checkHttpStatus(response);
  final decodedResponse =
      jsonDecode(utf8.decode(response.bodyBytes)) as List? ?? [];
  return decodedResponse.map((e) => Playlist.fromMap(e)).toList();
}

Future<List<Video>> listVideos({
  required String projectId,
  required UserCredentials creds,
}) async {
  final url = Uri.https(apiHost, 'video/list', {
    'p': projectId,
  });
  final response = await http.get(url, headers: {
    'AccessToken': creds.accessToken,
  });
  checkHttpStatus(response);
  final decodedResponse =
      jsonDecode(utf8.decode(response.bodyBytes)) as List? ?? [];
  return decodedResponse.map((e) => Video.fromMap(e)).toList();
}

Future<void> removeChannel({
  required String projectId,
  required String id,
  required bool blacklist,
  required UserCredentials creds,
}) async {
  final url = Uri.https(apiHost, 'channel/remove');
  final response = await http.post(url,
      headers: {'AccessToken': creds.accessToken},
      body: json.encode({
        'projectID': projectId,
        'id': id,
        'blacklist': blacklist,
      }));
  checkHttpStatus(response);
}

Future<void> removePlaylist({
  required String projectId,
  required String id,
  required bool blacklist,
  required UserCredentials creds,
}) async {
  final url = Uri.https(apiHost, 'playlist/remove');
  final response = await http.post(url,
      headers: {'AccessToken': creds.accessToken},
      body: json.encode({
        'projectID': projectId,
        'id': id,
        'blacklist': blacklist,
      }));
  checkHttpStatus(response);
}

Future<void> removeVideo({
  required String projectId,
  required String id,
  required bool blacklist,
  required UserCredentials creds,
}) async {
  final url = Uri.https(apiHost, 'video/remove');
  final response = await http.post(url,
      headers: {'AccessToken': creds.accessToken},
      body: json.encode({
        'projectID': projectId,
        'id': id,
        'blacklist': blacklist,
      }));
  checkHttpStatus(response);
}

Future<List<Marker>> getStack({
  required String projectId,
  required UserCredentials creds,
}) async {
  final url = Uri.https(apiHost, 'filter/stack', {
    'p': projectId,
  });
  final response = await http.get(url, headers: {
    'AccessToken': creds.accessToken,
  });
  checkHttpStatus(response);
  final decodedResponse = jsonDecode(utf8.decode(response.bodyBytes)) as Map;
  final markers = decodedResponse['markers'] as List;
  return markers.map((e) => Marker.fromMap(e)).toList();
}

Future<Marker> getRandomMarker({
  required String projectId,
  required UserCredentials creds,
}) async {
  final url = Uri.https(apiHost, 'randmarker', {
    'p': projectId,
  });
  final response = await http.get(url, headers: {
    'AccessToken': creds.accessToken,
  });
  checkHttpStatus(response);
  final decodedResponse = jsonDecode(utf8.decode(response.bodyBytes)) as Map;
  return Marker.fromMap(decodedResponse);
}

Future<void> classifyMarker({
  required String projectId,
  required String videoId,
  required int time,
  required bool label,
  required UserCredentials creds,
}) async {
  final url = Uri.https(apiHost, 'filter/classify');
  final response = await http.post(url,
      headers: {'AccessToken': creds.accessToken},
      body: json.encode({
        'projectID': projectId,
        'marker': {
          'videoID': videoId,
          'time': time,
        },
        'label': label ? 1 : 0,
      }));
  checkHttpStatus(response);
}

Future<bool> userExists(String username) async {
  final url = Uri.https(apiHost, 'user/exists', {
    'u': username,
  });
  final response = await http.get(url);
  if (response.statusCode != 200) {
    throw ErrorSummary("status ${response.statusCode}: ${response.body}");
  }
  final decodedResponse = jsonDecode(utf8.decode(response.bodyBytes)) as Map;
  return decodedResponse['exists'] as bool;
}

Future<bool> emailExists(String email) async {
  final url = Uri.https(apiHost, 'user/exists', {
    'e': email,
  });
  final response = await http.get(url);
  if (response.statusCode != 200) {
    throw ErrorSummary("status ${response.statusCode}: ${response.body}");
  }
  final decodedResponse = jsonDecode(utf8.decode(response.bodyBytes)) as Map;
  return decodedResponse['exists'] as bool;
}

Future<bool> projectExists(String name) async {
  final url = Uri.https(apiHost, 'project/exists', {
    'n': name,
  });
  final response = await http.get(url);
  if (response.statusCode != 200) {
    throw ErrorSummary("status ${response.statusCode}: ${response.body}");
  }
  final decodedResponse = jsonDecode(utf8.decode(response.bodyBytes)) as Map;
  return decodedResponse['exists'] as bool;
}

Future<UserCredentials> createAccount({
  required String username,
  required String email,
  required String firstName,
  required String lastName,
  required String password,
}) async {
  final url = Uri.https(apiHost, 'user/register');
  final response = await http.post(url,
      body: json.encode({
        'username': username,
        'email': email,
        'firstName': firstName,
        'lastName': lastName,
        'password': password,
      }));
  if (response.statusCode != 200) {
    throw ErrorSummary("status ${response.statusCode}: ${response.body}");
  }
  final decodedResponse = jsonDecode(utf8.decode(response.bodyBytes)) as Map;
  return UserCredentials(
    id: decodedResponse['id'],
    username: username,
    email: email,
    firstName: firstName,
    lastName: lastName,
    accessToken: decodedResponse['accessToken'],
    enabled: decodedResponse['enabled'],
  );
}

Future<void> signOut(UserCredentials creds) async {
  final url = Uri.https(apiHost, 'user/signout');
  final response = await http.post(url, headers: {
    'AccessToken': creds.accessToken,
  });
  if (response.statusCode != 200) {
    throw ErrorSummary("status ${response.statusCode}: ${response.body}");
  }
}

Future<UserCredentials> login(
  String username,
  String password,
) async {
  final url = Uri.https(apiHost, 'user/login');
  final response = await http.post(url,
      body: json.encode({
        'username': username,
        'password': password,
      }));
  if (response.statusCode == 401) {
    throw InvalidCredentialsError();
  } else if (response.statusCode != 200) {
    throw ErrorSummary("status ${response.statusCode}: ${response.body}");
  }
  final decodedResponse = jsonDecode(utf8.decode(response.bodyBytes)) as Map;
  return UserCredentials.fromMap(decodedResponse);
}

Future<List<Project>> listProjects({
  required UserCredentials creds,
}) async {
  final url = Uri.https(apiHost, 'project/list');
  final response = await http.get(url, headers: {
    'AccessToken': creds.accessToken,
  });
  checkHttpStatus(response);
  final decodedResponse =
      jsonDecode(utf8.decode(response.bodyBytes)) as List? ?? [];
  return decodedResponse.map((e) => Project.fromMap(e)).toList();
}

Future<Project> createProject({
  required UserCredentials creds,
  required String name,
}) async {
  final url = Uri.https(apiHost, 'project/create');
  final response = await http.post(url,
      headers: {
        'AccessToken': creds.accessToken,
      },
      body: json.encode({
        'name': name,
      }));
  checkHttpStatus(response);
  final decodedResponse = jsonDecode(utf8.decode(response.bodyBytes)) as Map;
  return Project.fromMap(decodedResponse);
}

Future<void> deleteProject({
  required UserCredentials creds,
  required String id,
}) async {
  final url = Uri.https(apiHost, 'project/delete');
  final response = await http.post(url,
      headers: {
        'AccessToken': creds.accessToken,
      },
      body: json.encode({
        'id': id,
      }));
  checkHttpStatus(response);
}

Future<List<SearchUser>> searchUsers({
  required UserCredentials creds,
  required String prefix,
}) async {
  final url = Uri.https(apiHost, 'user/search', {
    'p': prefix,
  });
  final response = await http.get(url, headers: {
    'AccessToken': creds.accessToken,
  });
  checkHttpStatus(response);
  final decodedResponse =
      jsonDecode(utf8.decode(response.bodyBytes)) as List? ?? [];
  return decodedResponse.map((e) => SearchUser.fromMap(e)).toList();
}

Future<void> addCollaborator({
  required UserCredentials creds,
  required String projectId,
  required String userId,
}) async {
  final url = Uri.https(apiHost, 'project/collaborators/add');
  final response = await http.post(url,
      headers: {
        'AccessToken': creds.accessToken,
      },
      body: json.encode({
        'userID': userId,
        'projectID': projectId,
      }));
  checkHttpStatus(response);
}

Future<void> removeCollaborator({
  required UserCredentials creds,
  required String projectId,
  required String userId,
}) async {
  final url = Uri.https(apiHost, 'project/collaborators/remove');
  final response = await http.post(url,
      headers: {
        'AccessToken': creds.accessToken,
      },
      body: json.encode({
        'userID': userId,
        'projectID': projectId,
      }));
  checkHttpStatus(response);
}

String videoThumbnail(String videoId) =>
    'https://$apiHost/video/thumbnail?v=$videoId';

String playlistThumbnail(String playlistId) =>
    'https://$apiHost/playlist/thumbnail?list=$playlistId';

String channelAvatar(String channelId) =>
    'https://$apiHost/channel/avatar?c=$channelId';
