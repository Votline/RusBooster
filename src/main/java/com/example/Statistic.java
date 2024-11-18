import java.sql.PreparedStatement;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.Statement;
import java.sql.ResultSet;

public class Statistic{
	private String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/statistic.db";
	private String sql;

	private String message;

	public void createTable(){
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "CREATE TABLE IF NOT EXISTS statistics ("+
				"user_id LONG NOT NULL UNIQUE," +
				"current_task INTEGER,"+
				"baddest_task STRING DEFAULT '1'," +
				"baddest_score INTEGER," +
				"better_task STRING DEFAULT '2'," +
				"better_score INTEGER," +
				"streak INTEGER" +
			");";
			Statement stmt = conn.createStatement();
			stmt.execute(sql);
		}
		catch(SQLException e){
			e.printStackTrace();
		}
	}
	public String getStatistic(String userName, long userId){
		createTable();
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "SELECT * FROM statistics WHERE user_id = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setLong(1, userId);
			ResultSet result = pstmt.executeQuery();
			if(result.next()){
				message = "👋Привет " + userName + "!" +
					"	\nТекущее задание: " + result.getString("current_task") +
					"	\nНаихудшая успеваимость: №" + result.getString("baddest_task") + ", " + result.getInt("baddest_score") +
					"	\nНаилучшая успеваимость: №" + result.getString("better_task") + ", " + result.getInt("better_score") +
					"\nТы занимаешься уже " + result.getInt("streak") + " дней подряд!👏";
			}
		}
		catch(SQLException e){
			e.printStackTrace();
			message = "Таблица " + "\"" + "statistics" + "\"" + " пуста для вашего userId";

		}
		System.out.println(message);
		return message;
	}
}
