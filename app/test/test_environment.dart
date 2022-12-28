import 'dart:convert';
import 'dart:io';
import 'package:flutter_test/flutter_test.dart';

class TestCredentials {
  final String firstName;
  final String lastName;
  final String email;
  final String username;
  final String password;

  TestCredentials({
    required this.firstName,
    required this.lastName,
    required this.email,
    required this.username,
    required this.password,
  });

  static TestCredentials fromMap(Map<dynamic, dynamic> m) => TestCredentials(
        firstName: requireJson<String>(m, 'firstName'),
        lastName: requireJson<String>(m, 'lastName'),
        email: requireJson<String>(m, 'email'),
        username: requireJson<String>(m, 'username'),
        password: requireJson<String>(m, 'password'),
      );
}

Type typeOf<T>() => T;

T requireJson<T>(Map<dynamic, dynamic> json, String key) {
  if (json[key] == null) {
    throw "missing required JSON field of type ${T.runtimeType}: $key";
  }
  if (json[key] is T) {
    return json[key] as T;
  } else {
    throw "error decoding required JSON field $key: expected type ${typeOf<T>()}";
  }
}

String extractPod(String service, List<String> pods) {
  late String line;
  try {
    line = pods.firstWhere((s) => s.contains(service));
  } catch (err) {
    expect(line, isNotNull,
        reason: "no pod listing for $service: ${err.toString()}");
  }
  expect(line.contains("Running"), isTrue, reason: "$service is not running");
  var pod = line.substring(line.indexOf(service));
  pod = pod.substring(0, pod.indexOf(" "));
  expect(pod, isNotEmpty, reason: "expected a pod name for $service");
  return pod;
}

class TestProject {
  final String name;
  final List<String> inputVideos;

  TestProject({
    required this.name,
    required this.inputVideos,
  });

  static TestProject fromMap(Map<dynamic, dynamic> m) => TestProject(
        name: m['name'],
        inputVideos: (m['inputVideos'] as List<dynamic>?)
                ?.map((e) => e as String)
                .toList() ??
            [],
      );
}

class TestEnvironment {
  late List<TestCredentials> creds;
  late String kubeContext;
  late String namespace;
  late String releaseName;
  late String gatewayPod;
  late String sourcesPod;
  late List<TestProject> projects;

  static const secretsJsonPath = 'test/secrets/test.json';

  static Future<TestEnvironment> create() async {
    final env = TestEnvironment();
    await env.setup();
    return env;
  }

  String get gatewayService => '$releaseName-gateway';
  String get sourcesService => '$releaseName-sources';

  Future<void> updatePods() async {
    final kubeListPods = await Process.run('kubectl', [
      'get',
      'pod',
      '-n',
      namespace,
      '--context',
      kubeContext,
    ]);
    expect(kubeListPods.exitCode, equals(0),
        reason:
            "kubectl get pod must succeed, make sure you have doctl authorization to our cluster, and kubectl (minimum client v1.18.6) installed: ${kubeListPods.stderr}");
    final List<String> pods =
        (kubeListPods.stdout as String).split("\n").sublist(1);
    gatewayPod = extractPod(gatewayService, pods);
    sourcesPod = extractPod(sourcesService, pods);
  }

  Future<void> exec(
    String podName,
    String command, [
    List<String> arguments = const [],
  ]) async {
    arguments = ['exec', '-n', namespace, podName, '--', command, ...arguments];
    command = 'kubectl';
    if (Platform.environment.containsKey("WSL")) {
      arguments = ['--exec', command, ...arguments];
      command = 'wsl.exe';
    }
    final cmd = await Process.run(command, arguments);
    expect(cmd.exitCode, equals(0),
        reason:
            "command \"$command ${arguments.join(' ')}\" must succeed, check k8s client compatibility");
  }

  Future<void> deleteTestUser() async {
    for (var cred in creds) {
      await exec(gatewayPod, 'gateway', [
        'iam',
        'delete-user',
        '--username',
        cred.username,
      ]);
    }
  }

  Future<void> deleteTestProject() async {
    for (var project in projects) {
      await exec(sourcesPod, 'sources', [
        'delete',
        'project',
        '--name',
        project.name,
      ]);
    }
  }

  Future<void> loadConfig() async {
    final json = jsonDecode(await File(secretsJsonPath).readAsString());
    kubeContext = requireJson<String>(json, 'kubeContext');
    namespace = requireJson<String>(json, 'namespace');
    releaseName = requireJson<String>(json, 'releaseName');
    projects =
        (json['projects'] as List).map((e) => TestProject.fromMap(e)).toList();
    creds =
        (json['creds'] as List).map((e) => TestCredentials.fromMap(e)).toList();
  }

  Future<void> setup() async {
    // there's an uncommitted json file containing
    // various secret values needed for the tests
    await loadConfig();

    // get the names of k8s pods for exec'ing
    await updatePods();

    // if the tests didn't tear down correctly
    // the last time, we'll have to delete the
    // test user and project here so we can
    // create them again
    await deleteTestUser();
    await deleteTestProject();
  }

  Future<void> tearDown() async {}
}
