import 'dart:convert';
import 'dart:io';
import 'dart:isolate';

import 'package:dart_scope_functions/dart_scope_functions.dart';
import 'package:meta/meta.dart';
import 'package:puby/api/api.exception.dart';
import 'package:puby/api/models/package.model.dart';
import 'package:puby/api/models/sdk_release.model.dart';
import 'package:puby/utils/types.dart';

abstract interface class IApiService {
  factory IApiService.instance() => _ApiService();

  Future<SDKReleaseModel> getSDKReleases();
  Future<List<PackageModel>> getLatestPackageVersions(List<String> packages);
}

@immutable
final class _ApiService implements IApiService {
  static const String _releasesHostname = 'storage.googleapis.com';
  static const String _pubAPIHostname = 'pub.dev';
  static const String _userAgent = 'puby:0.0.1 (+github.com/sunderee/puby)';

  final HttpClient _client;

  _ApiService() : _client = HttpClient();

  @override
  Future<SDKReleaseModel> getSDKReleases() {
    return _apiRequest<SDKReleaseModel>(
      _releasesHostname,
      'flutter_infra_release/releases/releases_macos.json',
      {
        HttpHeaders.acceptHeader: ContentType.json.toString(),
        HttpHeaders.contentTypeHeader: ContentType.json.toString(),
        HttpHeaders.userAgentHeader: _userAgent,
      },
      callback: (it) => SDKReleaseModel.fromJson(it),
    );
  }

  @override
  Future<List<PackageModel>> getLatestPackageVersions(
    List<String> packages,
  ) async {
    final List<Future<PackageModel>> futures = [];
    for (var item in packages) {
      final future = _apiRequest<PackageModel>(
        _pubAPIHostname,
        'api/packages/$item',
        {
          HttpHeaders.acceptHeader: 'application/vnd.pub.v2+json',
          HttpHeaders.contentTypeHeader: 'application/vnd.pub.v2+json',
          HttpHeaders.userAgentHeader: _userAgent,
        },
        callback: (it) => PackageModel.fromJson(it),
      );
      futures.add(future);
    }

    return Future.wait<PackageModel>(futures);
  }

  Future<T> _apiRequest<T>(
    String host,
    String endpoint,
    Headers headers, {
    required DecodingCallback<T> callback,
  }) async {
    return Isolate.run<T>(() async {
      final url = Uri.https(host, endpoint);
      final request = await _client.getUrl(url)
        ..headers.addAll(headers);
      final response = await request.close();

      if (response.statusCode >= 300) {
        throw ApiException(response);
      }
      final body = await response
          .transform(Utf8Decoder(allowMalformed: true))
          .reduce((previous, element) => previous + element);
      return (json.decode(body) as JsonObject).let((it) => callback.call(it));
    });
  }
}

extension _HttpHeadersExt on HttpHeaders {
  void addAll(Headers headers) => headers.forEach((k, v) => add(k, v));
}
