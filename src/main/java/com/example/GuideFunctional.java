package com.example;

import java.sql.PreparedStatement;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.Statement;

import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.InlineKeyboardMarkup;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.buttons.KeyboardButton;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.buttons.InlineKeyboardButton;

public class GuideFunctional{
	private void createTable(){
		String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/guide.db";
		String sql;
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "CREATE TABLE IF NOT EXISTS guides (" +
				"taskId INTEGER PRIMARY KEY," +
				"guide TEXT NOT NULL DEFAULT '0'" +
			");";
			Statement stmt = conn.createStatement();
			stmt.execute(sql);
			stmt.close();
		}
		catch(SQLException e){
			e.printStackTrace();
		}
	}

	private int getUserTaskId(long userId){
		String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/guide.db";
		String sql;
		int taskId = 0;
		url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/statistic.db";
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "SELECT current_task FROM statistics WHERE user_id = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setLong(1, userId);
			ResultSet result = pstmt.executeQuery();
			taskId = result.getInt("current_task");
			result.close(); pstmt.close();
		}
		catch(SQLException e){
			createTable();
			e.printStackTrace();
		}
		return taskId;
	}

	public SendMessage sendGuide(long userId){
		String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/guide.db";
		String sql;
		SendMessage guideMessage = new SendMessage();
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "SELECT guide FROM guides WHERE taskId = ?";
			PreparedStatement sendPstmt = conn.prepareStatement(sql);
			sendPstmt.setInt(1, getUserTaskId(userId));
			ResultSet result = sendPstmt.executeQuery();
			String message = "Error";
			if(result.next()){
				message = result.getString("guide");
			}
			guideMessage.setChatId(String.valueOf(userId));
			guideMessage.setText(message);
			
			sendPstmt.close(); result.close();
		}
		catch(SQLException e){
			createTable();
			e.printStackTrace();
		}
		return guideMessage;
	}

	public void insertGuide(int taskId, String guide){
		String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/guide.db";
		String sql;
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "INSERT INTO guides (taskId, guide) VALUES (?, ?)";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setInt(1, taskId);
			pstmt.setString(2, guide);
			pstmt.executeUpdate();
			pstmt.close();
		}
		catch(SQLException e){
			e.printStackTrace();
		}
	}

	public void removeGuide(int taskId){
		String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/guide.db";
		String sql;
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "DELETE FROM guides WHERE taskId = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setInt(1, taskId);
			pstmt.executeUpdate();
			pstmt.close();
		}
		catch(SQLException e){
			e.printStackTrace();
		}
	}
}
