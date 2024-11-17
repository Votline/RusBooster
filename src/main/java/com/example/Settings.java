import java.sql.PreparedStatement;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.Statement;

public class Settings{
	private String message;

	private String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/statistic.db";
	private String sql;

	private Statistic statistic = new Statistic();

	public void chooseExercise(long userId, Integer numberOfExercise){
		try(Connection conn = DriverManager.getConnection(url)){
			statistic.createTable();
			sql = "INSERT OR REPLACE INTO statistics (user_id, current_task) VALUES (?, ?)";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setLong(1, userId);
			pstmt.setInt(2, numberOfExercise);
			pstmt.executeUpdate();
			message = "Текущее упражнение для пользователя " + userId + ": " + numberOfExercise;
			System.out.println(message);
		}
		catch(SQLException e){
			e.printStackTrace();
		}
	}
}
