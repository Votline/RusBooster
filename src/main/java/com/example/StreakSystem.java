import java.time.LocalDate;

import java.sql.PreparedStatement;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.ResultSet;


public class StreakSystem{
	private String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/statistic.db";
	private String sql;
	public void checkStreak(long userId){
		UserState userState = UserStateManager.getUserState(userId);
		int lastActiveDate = 0;
		int currentDate = (int) LocalDate.now().toEpochDay();
		int streak = 0;

		try(Connection conn = DriverManager.getConnection(url)){
			sql = "SELECT last_Active_Date, streak FROM statistics WHERE user_id = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setLong(1, userId);
			ResultSet result = pstmt.executeQuery();
			while(result.next()){
				lastActiveDate = result.getInt("last_Active_Date");	
				streak = result.getInt("streak");	
			}
			if(currentDate != lastActiveDate){
				userState.isActive = true;
				streak += 1;
			}
			sql = "UPDATE statistics SET last_Active_Date = ?, streak = ? WHERE user_id = ?";
			pstmt = conn.prepareStatement(sql);
			pstmt.setInt(1, lastActiveDate);
			pstmt.setInt(2, streak);
			pstmt.setLong(3, userId);
			pstmt.executeUpdate();
		}
		catch(SQLException e){
			e.printStackTrace();
		}
		
	}
}
