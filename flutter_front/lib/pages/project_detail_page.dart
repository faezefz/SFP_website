import 'dart:convert';
import 'dart:typed_data';
import 'package:flutter/material.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter_project/pages/colors.dart';
import 'package:flutter_project/pages/results_page.dart';
import 'package:http/http.dart';
import 'api_service.dart';

class ProjectDetailsPage extends StatefulWidget {
  final String projectId;
  final int userId;

  const ProjectDetailsPage(
      {super.key, required this.projectId, required this.userId});

  @override
  _ProjectDetailsPageState createState() => _ProjectDetailsPageState();
}

class _ProjectDetailsPageState extends State<ProjectDetailsPage> {
  bool isLoading = true;
  Map<String, dynamic> projectDetails = {};
  List<String> models = [];
  String selectedModel = '';
  List<String> allDatasets = [];
  List<String> selectedTrainDatasets = [];
  List<String> selectedTestDatasets = [];
  bool isCustomModelSelected = false;
  Uint8List? customModel;

  @override
  void initState() {
    super.initState();
    _fetchProjectDetails();
  }

  Future<void> _fetchProjectDetails() async {
    try {
      int projectId = int.parse(widget.projectId);
      Map<String, dynamic> details =
          await ApiService().getProjectById(projectId);

      List<String> datasets =
          await ApiService().getDatasetsByUserId(widget.userId);

      setState(() {
        projectDetails = details;
        models = List<String>.from(details['models'] ?? []);
        allDatasets = datasets;
        isLoading = false;
      });
    } catch (e) {
      setState(() {
        isLoading = false;
      });
    }
  }

  Future<void> _uploadDataset(String type) async {
    FilePickerResult? result = await FilePicker.platform.pickFiles(
      type: FileType.custom,
      allowedExtensions: ['csv'],
    );

    if (result != null) {
      Uint8List fileBytes = result.files.single.bytes!;
      String datasetName = '';
      String datasetDescription = '';

      await _showDatasetDialog((name, description) {
        datasetName = name;
        datasetDescription = description;
      });

      if (datasetName.isNotEmpty && datasetDescription.isNotEmpty) {
        setState(() {
          isLoading = true;
        });

        bool success = await ApiService().importDataset(
          widget.userId,
          datasetName,
          datasetDescription,
          fileBytes,
        );

        setState(() {
          isLoading = false;
        });

        if (success) {
          setState(() {
            if (type == 'Train') {
              selectedTrainDatasets.add(datasetName);
            } else {
              selectedTestDatasets.add(datasetName);
            }
          });

          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('$type Dataset uploaded successfully!')),
          );
        } else {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Failed to upload $type dataset')),
          );
        }
      }
    }
  }

  Future<void> _showDatasetDialog(Function(String, String) onSubmit) async {
    String name = '';
    String description = '';

    await showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          backgroundColor: AppColors.darktGry,
          title: const Text('Enter Dataset Name and Description',
              style: TextStyle(color: AppColors.lightGr)),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                onChanged: (value) {
                  name = value;
                },
                decoration: InputDecoration(
                  hintText: 'Dataset Name',
                  hintStyle: TextStyle(color: AppColors.lightGry),
                  focusedBorder: UnderlineInputBorder(
                    borderSide: BorderSide(color: AppColors.lightOr),
                  ),
                ),
              ),
              const SizedBox(height: 8),
              TextField(
                onChanged: (value) {
                  description = value;
                },
                decoration: InputDecoration(
                  hintText: 'Dataset Description',
                  hintStyle: TextStyle(color: AppColors.lightGry),
                  focusedBorder: UnderlineInputBorder(
                    borderSide: BorderSide(color: AppColors.lightOr),
                  ),
                ),
              ),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () {
                Navigator.of(context).pop();
              },
              child: const Text('Cancel',
                  style: TextStyle(color: AppColors.lightGr)),
            ),
            TextButton(
              onPressed: () {
                onSubmit(name, description);
                Navigator.of(context).pop();
              },
              child: const Text('Submit',
                  style: TextStyle(color: AppColors.lightGr)),
            ),
          ],
        );
      },
    );
  }

  void _removeDataset(String dataset, String type) {
    setState(() {
      if (type == 'Train') {
        selectedTrainDatasets.remove(dataset);
      } else {
        selectedTestDatasets.remove(dataset);
      }
    });
  }

  Future<void> _startPrediction() async {
    if (selectedTrainDatasets.isEmpty ||
        selectedTestDatasets.isEmpty ||
        selectedModel.isEmpty) {
      return;
    }

    Uint8List? modelToSend = isCustomModelSelected ? customModel : null;

    Response response = await ApiService().startPrediction(
      int.parse(widget.projectId),
      selectedTrainDatasets,
      selectedTestDatasets,
      selectedModel,
      modelToSend,
      isCustomModelSelected,
    );

    if (response.statusCode == 200) {
      final responseData = json.decode(response.body);

      List<Map<String, dynamic>> result =
          List<Map<String, dynamic>>.from(responseData['result']);

      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Prediction started successfully')),
      );

      Navigator.push(
        context,
        MaterialPageRoute(
          builder: (context) => ResultsPage(result: result),
        ),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Failed to start prediction')),
      );
    }
  }

  Widget _buildModelSelection() {
    List<String> machineLearningModels = [
      'Linear Regression',
      'Decision Tree',
      'Random Forest',
      'SVM',
      'KNN',
      'Neural Network',
      'Logistic Regression',
      'Naive Bayes',
      'XGBoost'
    ];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('Choose Model:', style: TextStyle(color: AppColors.lightGr)),
        DropdownButton<String>(
          value: selectedModel.isEmpty ? null : selectedModel,
          hint: const Text('Select Model'),
          onChanged: (value) {
            setState(() {
              selectedModel = value!;
              isCustomModelSelected = false;
            });
          },
          items: machineLearningModels
              .map<DropdownMenuItem<String>>((String model) {
            return DropdownMenuItem<String>(
              value: model,
              child: Text(model, style: TextStyle(color: AppColors.darkGr)),
            );
          }).toList(),
        ),
        CheckboxListTile(
          title: const Text('Use Custom Model',
              style: TextStyle(color: AppColors.lightGr)),
          value: isCustomModelSelected,
          onChanged: (value) {
            setState(() {
              isCustomModelSelected = value!;
            });
          },
        ),
        if (isCustomModelSelected)
          ElevatedButton(
            onPressed: _selectCustomModel,
            style: ElevatedButton.styleFrom(primary: AppColors.lightOr),
            child: const Text('Upload Custom Model'),
          ),
      ],
    );
  }

  Future<void> _selectCustomModel() async {
    Uint8List? customModelFile = await _uploadCustomModel();

    setState(() {
      customModel = customModelFile;
      isCustomModelSelected = true;
      selectedModel = "";
    });
  }

  Future<Uint8List> _uploadCustomModel() async {
    FilePickerResult? result = await FilePicker.platform.pickFiles(
      type: FileType.custom,
      allowedExtensions: ['pkl'],
    );

    if (result != null) {
      Uint8List fileBytes = result.files.single.bytes!;
      return fileBytes;
    }
    return Uint8List(0);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Project Details',
            style: TextStyle(color: AppColors.lightGr)),
        backgroundColor: AppColors.darkGr,
      ),
      body: isLoading
          ? const Center(child: CircularProgressIndicator())
          : Padding(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                children: [
                  Text(
                    'Project Name: ${projectDetails['name']}',
                    style: TextStyle(color: AppColors.darkGr),
                  ),
                  Text(
                    'Description: ${projectDetails['description'] ?? 'No description!'}',
                    style: TextStyle(color: AppColors.darkGr),
                  ),
                  const SizedBox(height: 16),
                  _buildDatasetSelection(selectedTrainDatasets, 'Train'),
                  _buildDatasetSelection(selectedTestDatasets, 'Test'),
                  _buildModelSelection(),
                  const SizedBox(height: 16),
                  ElevatedButton(
                    onPressed: _startPrediction,
                    style: ElevatedButton.styleFrom(primary: AppColors.lightOr),
                    child: const Text('Start Prediction'),
                  ),
                ],
              ),
            ),
    );
  }

  Widget _buildDatasetSelection(List<String> datasets, String type) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('$type Dataset(s):', style: TextStyle(color: AppColors.lightGr)),
        ...datasets.map((dataset) {
          return Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(dataset, style: TextStyle(color: AppColors.darkGr)),
              IconButton(
                icon: const Icon(Icons.delete),
                onPressed: () {
                  _removeDataset(dataset, type);
                },
                color: AppColors.lightOr,
              ),
            ],
          );
        }).toList(),
        Row(
          children: [
            Expanded(
              child: DropdownButton<String>(
                value: datasets.isNotEmpty ? datasets[0] : null,
                hint: const Text('Import from your datasets'),
                onChanged: (value) {
                  setState(() {
                    if (type == 'Train') {
                      if (!selectedTrainDatasets.contains(value)) {
                        selectedTrainDatasets.add(value!);
                      }
                    } else {
                      if (!selectedTestDatasets.contains(value)) {
                        selectedTestDatasets.add(value!);
                      }
                    }
                  });
                },
                items:
                    allDatasets.map<DropdownMenuItem<String>>((String dataset) {
                  return DropdownMenuItem<String>(
                    value: dataset,
                    child: Text(dataset,
                        style: TextStyle(color: AppColors.darkGr)),
                  );
                }).toList(),
              ),
            ),
            ElevatedButton(
              onPressed: () {
                _uploadDataset(type);
              },
              style: ElevatedButton.styleFrom(primary: AppColors.lightOr),
              child: const Text('Upload'),
            ),
          ],
        ),
      ],
    );
  }
}
