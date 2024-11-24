import java.sql.PreparedStatement;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.ResultSet;

import java.util.HashMap;

public class MainFunctional{
	private String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/words.db";
	private String urlStat = "jdbc:sqlite:" + System.getProperty("user.dir") + "/statistic.db";
	private String sql;
	public void makeTask(long userId){
		int task_Id = 0;
		try(Connection connSet = DriverManager.getConnection(urlStat)){
			sql = "SELECT current_task FROM statistics WHERE user_id = ?";
			PreparedStatement pstmt = connSet.prepareStatement(sql);
			pstmt.setString(1, userId);
			task_Id = pstmt.executeQuery().next().getInteger("current_task");

		}
		catch(SQLException e){
			e.printStackTrace();
		}
		try(Connection conn = DriverManager.getConnection(url)){
			HashMap<String, String> task = new HashMap<>();
			sql = "SELECT word, explanation FROM words WHERE task_id = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setString(1, task_Id);
			ResultSet result = pstmt.executeQuery();

			while(result.next()){
				String word = result.getString("word");
				String explanation = result.getString("explanation");
				task.put(word, explanation);
				System.out.println(word);
				System.out.println(explanation);
			}
		}
		catch(SQLException e){
			e.printStackTrace();
		}
	}

}
