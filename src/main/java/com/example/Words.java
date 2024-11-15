import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.Statement;
import java.sql.PreparedStatement;

import java.io.File;
import java.io.IOException;
import java.util.Map;
import java.util.HashMap;

public class Words{
	private Map<String, String> words = new HashMap<>();
	private String message;

	private String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/words.db";
	private String sql;

	private void createTable(){
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "CREATE TABLE IF NOT EXISTS words (" +
				"word TEXT NOT NULL," +
				"explanation TEXT NOT NULL" +
			");";
			Statement statement = conn.createStatement();
			statement.execute(sql);
			System.out.println("Создал таблицу");
		}
		catch(SQLException e){
			e.printStackTrace();
		}

	}
	public String addWord(String name, String explanation){
		try(Connection conn = DriverManager.getConnection(url)){
			createTable();
			sql = "INSERT INTO words(word, explanation) VALUES(?, ?)";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setString(1, name);
			pstmt.setString(2, explanation);
			pstmt.executeUpdate();
			message = "Добавил слово " + name + ". Значение: " + explanation;
			System.out.println(message);
		}
		catch(SQLException e){
			e.printStackTrace();
		}
		return message;

	}
/*
	public String removeWord(String name){
		try{
			if(base.exists()){
				if(words.containsKey(name)){
					message = "Удалил слово " + name;
					System.out.println(message);
				}
				else{
					message = "Слово" + name + " не найдено! ";
					System.out.println(message);
				}
			}
		}
		catch(IOException e){
			e.printStackTrace();

		}
		return message;
	}

	public String showAllBase(){
		try{
			if(base.exists()){
				message = " ";
				words = objMapper.readValue(base, Map.class);
				for(Map.Entry<String, String> entry : words.entrySet()){
					message += "\nСлово: " + entry.getKey() + "\nЗначение: " + entry.getValue();
					System.out.println(message);
				}
			}
		}
		catch(IOException e){
			e.printStackTrace();
		}
		return message;
	}
*/
}
