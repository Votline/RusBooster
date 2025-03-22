package com.example;

import java.sql.PreparedStatement;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.Statement;
import java.sql.ResultSet;

import java.util.ArrayList;
import java.util.List;

public class Words{
	private String message;

	private String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/words.db";
	private String sql;

	private void createTable(){
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "CREATE TABLE IF NOT EXISTS words (" +
				"word TEXT NOT NULL," +
				"explanation TEXT NOT NULL," +
				"task_id INTEGER NOT NULL" +
			");";
			Statement stmt = conn.createStatement();
			stmt.execute(sql);
		}
		catch(SQLException e){
			e.printStackTrace();
		}

	}
	public String addWord(String name, String explanation, int task_id){
		createTable();
		try(Connection conn = DriverManager.getConnection(url)){
			if(checkBase(name, task_id)){
				message = "Слово " + name + " для задания " + task_id + " уже существует в базе данных.";
				System.out.println(message);
				return message;
			}

			sql = "INSERT INTO words(word, explanation, task_id) VALUES(?, ?, ?)";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setString(1, name);
			pstmt.setString(2, explanation);
			pstmt.setInt(3, task_id);
			pstmt.executeUpdate();
			message = "Добавил слово " + name + ". Значение: " + explanation + ". Задание: " + task_id + ".";
		}
		catch(SQLException e){
			e.printStackTrace();
		}
		return message;

	}

	public String removeWord(String name, int task_id){
		try(Connection conn = DriverManager.getConnection(url)){
			if( !(checkBase(name, task_id)) ){
				message = "Слово " + name + " для задания " + task_id + " не существует в базе данных.";
				System.out.println(message);
				return message;
			}
			sql = "DELETE FROM words WHERE word = ? AND task_id = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setString(1, name);
			pstmt.setInt(2, task_id);
			pstmt.executeUpdate();
			message = "Удалил слово " + name + " для задания " + task_id + " из базы данных.";
		}
		catch(SQLException e){
			e.printStackTrace();

		}
		return message;
	}

	public List<String> showAllBase(int task_id, String whatNeed){
		List<String> allWords = new ArrayList<>(); int cnt = 0;
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "SELECT word, explanation FROM words WHERE task_id = ?";
			if(!"null".equals(whatNeed)){
				sql += "AND explanation COLLATE utf8mb4_bin LIKE ?";	
			}
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setInt(1, task_id);
			if(!"null".equals(whatNeed)){
				pstmt.setString(2, "%" + whatNeed + "%");
			}
			ResultSet result = pstmt.executeQuery();

			StringBuilder message = new StringBuilder();
			while(result.next()){
				String word = result.getString("word");
				String explanation = result.getString("explanation");
				
				if(!"null".equals(whatNeed) && whatNeed.length() <= 1){
					String potentialWord = explanation.split(" ")[0];
					if(!potentialWord.contains(whatNeed)){
						continue;
					}
				}

				cnt++;
				message.append("Слово: \"")
					.append(word)
					.append("\". Значение: \"")
					.append(explanation)
					.append("\". Номер задания: \"")
					.append(task_id)
					.append("\"\n\n");
				if(cnt == 5){
					allWords.add(message.toString());
					message.setLength(0);
					cnt = 0;
				}
			}
			if(cnt != 0) {allWords.add(message.toString());}
			if(allWords.isEmpty()){allWords.add("Ничего не найдено");}
		}
		catch(SQLException e){
			e.printStackTrace();
		}
		return allWords;
	}

	private boolean checkBase(String name, int task_id){
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "SELECT 1 FROM words WHERE word = ? AND task_id = ?";
			PreparedStatement checkPSTMT = conn.prepareStatement(sql);
			checkPSTMT.setString(1, name);
			checkPSTMT.setInt(2, task_id);
			if( (checkPSTMT.executeQuery()).next() ){
				return true;
			}
		}
		catch(SQLException e){
			e.printStackTrace();
		}
		return false;
	}
}
