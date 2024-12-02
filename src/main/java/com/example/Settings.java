import org.telegram.telegrambots.meta.api.methods.send.SendMessage;

import java.sql.PreparedStatement;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.Statement;

public class Settings{
	private String messageText;

	private String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/statistic.db";
	private String sql;

	private Statistic statistic = new Statistic();

	public void chooseExercise(long userId, String message){
		int numberOfExercise = Integer.parseInt(message);
		if(numberOfExercise <= 26 && numberOfExercise >= 1){
			try(Connection conn = DriverManager.getConnection(url)){
				statistic.createTable();
				sql = "INSERT OR REPLACE INTO statistics (user_id, current_task) VALUES (?, ?)";
				PreparedStatement pstmt = conn.prepareStatement(sql);
				pstmt.setLong(1, userId);
				pstmt.setInt(2, numberOfExercise);
				pstmt.executeUpdate();			
			}
			catch(SQLException e){
				e.printStackTrace();
			}
		}
	}
}
