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
  final List<VideoListItem> videos;
  final DateTime timestamp;

  Dataset({
    required this.id,
    required this.videos,
    required this.timestamp,
  });

  static Dataset fromMap(Map m) => Dataset(
        id: m['id'],
        timestamp: DateTime.fromMicrosecondsSinceEpoch(m['timestamp'] ~/ 1000),
        videos:
            (m['videos'] as List).map((o) => VideoListItem.fromMap(o)).toList(),
      );
}

class VideoListItem {
  final String id;
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
  bool blacklist;

  VideoListItem({
    required this.uploader,
    required this.uploaderId,
    required this.thumbnail,
    required this.id,
    required this.title,
    required this.channel,
    required this.channelId,
    required this.duration,
    required this.width,
    required this.height,
    required this.uploadDate,
    required this.fps,
    this.blacklist = false,
  });

  static VideoListItem fromMap(Map m) => VideoListItem(
        id: m['id'],
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
        blacklist: m['blacklist'] ?? false,
      );
}

class PlaylistListItem {
  final String id;
  final String title;
  final String channel;
  final String channelId;
  final int numVideos;

  bool blacklist;

  PlaylistListItem({
    required this.id,
    required this.title,
    required this.channel,
    required this.channelId,
    required this.numVideos,
    this.blacklist = false,
  });

  static PlaylistListItem fromMap(Map m) => PlaylistListItem(
        id: m['id'],
        title: m['title'],
        channel: m['channel'],
        channelId: m['channelID'],
        numVideos: m['numVideos'],
        blacklist: m['blacklist'],
      );
}

class ChannelListItem {
  final String id;
  final String name;
  final String avatarUrl;
  bool blacklist;

  ChannelListItem({
    required this.id,
    required this.name,
    required this.avatarUrl,
    this.blacklist = false,
  });

  static ChannelListItem fromMap(Map m) => ChannelListItem(
        id: m['id'],
        name: m['name'],
        avatarUrl: m['avatar'],
        blacklist: m['blacklist'],
      );
}

class SourceInput {
  List<VideoListItem> videos = [];
  List<PlaylistListItem> playlists = [];
  List<ChannelListItem> channels = [];
}

class SourceOutput {
  List<VideoListItem> videos = [];
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

Future<ChannelListItem?> addChannel({
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
  return ChannelListItem.fromMap(decodedResponse);
}

Future<PlaylistListItem?> addPlaylist({
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
  return PlaylistListItem.fromMap(decodedResponse);
}

Future<VideoListItem?> addVideo({
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
  return VideoListItem.fromMap(decodedResponse);
}

Future<List<ChannelListItem>> listChannels({
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
  return decodedResponse.map((e) => ChannelListItem.fromMap(e)).toList();
}

Future<List<PlaylistListItem>> listPlaylists({
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
  return decodedResponse.map((e) => PlaylistListItem.fromMap(e)).toList();
}

Future<List<VideoListItem>> listVideos({
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
  return decodedResponse.map((e) => VideoListItem.fromMap(e)).toList();
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
