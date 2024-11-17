import java.sql.PreparedStatement;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.Statement;
import java.sql.ResultSet;

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
			System.out.println("Создал таблицу words ");
		}
		catch(SQLException e){
			e.printStackTrace();
		}

	}
	public String addWord(String name, String explanation, String task_id){
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
			pstmt.setString(3, task_id);
			pstmt.executeUpdate();
			message = "Добавил слово " + name + ". Значение: " + explanation + ". Задание: " + task_id + ".";
			System.out.println(message);
		}
		catch(SQLException e){
			e.printStackTrace();
		}
		return message;

	}

	public String removeWord(String name, String task_id){
		try(Connection conn = DriverManager.getConnection(url)){
			if( !(checkBase(name, task_id)) ){
				message = "Слово " + name + " для задания " + task_id + " не существует в базе данных.";
				System.out.println(message);
				return message;
			}
			sql = "DELETE FROM words WHERE word = ? AND task_id = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setString(1, name);
			pstmt.setString(2, task_id);
			pstmt.executeUpdate();
			message = "Удалил слово " + name + " для задания " + task_id + " из базы данных.";
			System.out.println(message);
		}
		catch(SQLException e){
			e.printStackTrace();

		}
		return message;
	}

	public String showAllBase(){
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "SELECT word, explanation, task_id FROM words";
			Statement stmt = conn.createStatement();
			ResultSet result = stmt.executeQuery(sql);
			message = "";
			while(result.next()){
				String word = result.getString("word");
				String explanation = result.getString("explanation");
				String task_id = result.getString("task_id");
				message += "Слово: " + "\"" + word + "\"" + ". Значение: " + "\"" +  explanation + "\"" +  ". Номер задания: " + "\"" +  task_id + "\"" + "\n \n";
				System.out.println(message);
			}
		}
		catch(SQLException e){
			e.printStackTrace();
		}
		return message;
	}

	private boolean checkBase(String name, String task_id){
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "SELECT 1 FROM words WHERE word = ? AND task_id = ?";
			PreparedStatement checkPSTMT = conn.prepareStatement(sql);
			checkPSTMT.setString(1, name);
			checkPSTMT.setString(2, task_id);
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
