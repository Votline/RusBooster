import java.sql.PreparedStatement;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.Statement;

public class Settings{
	private String message;

	private String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/settings.db";
	private String sql;

	private void createTable(){
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "CREATE TABLE IF NOT EXISTS settings (" +
				"user_id STRING NOT NULL UNIQUE," +
				"task_id STRING NOT NULL" +
			");";
			Statement statement = conn.createStatement();
			statement.execute(sql);
			System.out.println("Создал таблицу settings ");
		}
		catch(SQLException e){
			e.printStackTrace();
		}
	}

	public void chooseExercise(String userId, String numberOfExercise){
		try(Connection conn = DriverManager.getConnection(url)){
			createTable();
			sql = "INSERT OR REPLACE INTO settings (user_id, task_id) VALUES (?, ?)";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setString(1, userId);
			pstmt.setString(2, numberOfExercise);
			pstmt.executeUpdate();
			message = "Текущее упражнение для пользователя " + userId + " :" + numberOfExercise;
			System.out.println(message);
		}
		catch(SQLException e){
			e.printStackTrace();
		}
	}
}
