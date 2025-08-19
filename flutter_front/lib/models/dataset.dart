import 'dart:typed_data';

class Dataset {
  final int id;
  final int userId;
  final String name;
  final String description;
  final Uint8List content;
  final String uploadedAt;

  Dataset({
    required this.id,
    required this.userId,
    required this.name,
    required this.description,
    required this.content,
    required this.uploadedAt,
  });

  factory Dataset.fromJson(Map<String, dynamic> json) {
    return Dataset(
      id: json['id'],
      userId: json['user_id'],
      name: json['name'],
      description: json['description'],
      content: json['content'] != null
          ? Uint8List.fromList(List<int>.from(json['content']))
          : Uint8List(0),
      uploadedAt: json['uploaded_at'],
    );
  }
}
