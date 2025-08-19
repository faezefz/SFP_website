class Project {
  final int id;
  final int ownerUserId;
  final String name;
  final String description;
  final String visibility;
  final String createdAt;

  Project({
    required this.id,
    required this.ownerUserId,
    required this.name,
    required this.description,
    required this.visibility,
    required this.createdAt,
  });

  factory Project.fromJson(Map<String, dynamic> json) {
    return Project(
      id: json['id'],
      ownerUserId: json['owner_user_id'],
      name: json['name'],
      description: json['description'],
      visibility: json['visibility'],
      createdAt: json['created_at'],
    );
  }
}
