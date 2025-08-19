import 'package:flutter/material.dart';
import 'colors.dart';

class ResultsPage extends StatelessWidget {
  final List<Map<String, dynamic>> result;

  ResultsPage({required this.result});

  @override
  Widget build(BuildContext context) {
    if (result.isEmpty) {
      return Scaffold(
        appBar: AppBar(
          title: Text('Prediction Results',
              style: TextStyle(color: AppColors.lightGr)),
          backgroundColor: AppColors.darkGr,
        ),
        body: Center(
          child: Text('No data available',
              style: TextStyle(color: AppColors.darkGr)),
        ),
      );
    }

    List<String> columns = result[0].keys.toList();

    return Scaffold(
      appBar: AppBar(
        title: Text('Prediction Results',
            style: TextStyle(color: AppColors.lightGr)),
        backgroundColor: AppColors.darkGr,
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Prediction Results:',
              style: TextStyle(color: AppColors.lightGr, fontSize: 22),
            ),
            SizedBox(height: 20),
            Expanded(
              child: SingleChildScrollView(
                scrollDirection: Axis.horizontal,
                child: DataTable(
                  columns: columns.map((col) {
                    return DataColumn(
                        label: Text(col,
                            style: TextStyle(color: AppColors.lightGr)));
                  }).toList(),
                  rows: result.map((row) {
                    return DataRow(
                      cells: columns.map((col) {
                        return DataCell(Text(row[col].toString(),
                            style: TextStyle(color: AppColors.darkGr)));
                      }).toList(),
                    );
                  }).toList(),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
