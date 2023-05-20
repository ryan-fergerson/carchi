# Carchi, Chat archive

Save your ChatGPT conversations for backup and to do your own personal data analytics.

Search through your conversations using PostgreSQL's full-text searching features.

## Overview

This program imports conversation data from a provided ZIP file into a database. The input data is structured in JSON format based on the data export from the ChatGPT web application.

## Structure of Input Data

The program expects the input data to have a specific structure. Here is the definition of the conversation data:

- A `Conversation` has an
  - `Id`, `Title`, `CreateTime`, `UpdateTime`, `Mapping`, `ModerationResults`, `CurrentNode`, and `PluginIds`.
- The `Mapping` field contains a map of `Node` objects.
- Each `Node` has an
  - `Id`, an optional `Message`, an optional `Parent`, and a list of `Children` node ids.
- A `NodeMessage` has a
  - `Id`, `Author`, `CreateTime`, `UpdateTime`, `Content`, `EndTurn`, `Weight`, `Metadata` and `Recipient`.
- The `Content` of a `NodeMessage` has a `ContentType` and a list of `Parts`.

## How the Program Works

1. The program starts by reading the input from a ZIP file provided as a command-line argument. The ZIP file should contain a JSON file named "conversations.json".

2. The program then parses the input data into a slice of `Conversation` objects.

3. A connection to the database is opened.

4. The program starts a database transaction and prepares SQL statements for inserting conversations, nodes, and messages.

5. For each `Conversation`, the program executes the conversation statement to insert a row into the conversation table in the database.

6. For each `Node` in the `Mapping` field of the `Conversation`, the program executes the node statement to insert a row into the node table. If the `Node` has a `Message`, the program also executes the message statement to insert a row into the message table.

7. Finally, the transaction is committed to the database.

## Error Handling

The program uses custom error types to provide detailed information about any errors that occur during execution. For example, an error might occur when reading the ZIP file, parsing the JSON data, connecting to the database, preparing the SQL statements, executing the SQL statements, or committing the transaction.

Each error is associated with an action, such as "reading zip file" or "executing conversation statement". If an error occurs, the program logs the error and exits.
