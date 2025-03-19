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

import java.util.Collections;
import java.util.Arrays;
import java.util.List;

public class GuideFunctional{
	private void createTable(){
		String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/guide.db";
		String sql;
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "CREATE TABLE IF NOT EXISTS guides (" +
				"task_id INTEGER PRIMARY KEY," +
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
			sql = "SELECT guide FROM guides WHERE task_id = ?";
			PreparedStatement sendPstmt = conn.prepareStatement(sql);
			sendPstmt.setInt(1, getUserTaskId(userId));
			ResultSet result = sendPstmt.executeQuery();
			String message = "Error";
			if(result.next()){
				InlineKeyboardMarkup guideKeyboard = new InlineKeyboardMarkup();
					InlineKeyboardButton back = new InlineKeyboardButton();
					InlineKeyboardButton next = new InlineKeyboardButton();
					InlineKeyboardButton allPages = new InlineKeyboardButton();
					back.setCallbackData("back"); next.setCallbackData("next"); allPages.setCallbackData("toMain");
				List<InlineKeyboardButton> row = Arrays.asList(back, next);
				guideKeyboard.setKeyboard(Arrays.asList(
					row, 
					Collections.singletonList(allPages)
				));
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
			sql = "INSERT INTO guides (task_id, guide) VALUES (?, ?)";
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
			sql = "DELETE FROM guides WHERE task_id = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setInt(1, taskId);
			pstmt.executeUpdate();
			pstmt.close();
		}
		catch(SQLException e){
			e.printStackTrace();
		}
	}
	private StringBuilder showAllBase(int taskId){
		String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/guide.db";
		String sql;

		StringBuilder message = new StringBuilder();
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "SELECT guide FROM guides WHERE task_id = ?";
			
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setInt(0, taskId);
			ResultSet result = pstmt.executeQuery();
			
			if(result.next()){
				String guide = result.getString("guide");
				String[] themes = guide.split("!");
				for(String theme : themes){
					if(theme.length() > 511){
						int maxLenght = 509;
						for(int i = 0; i < theme.length(); i += maxLenght){
							String part = theme.substring(i, Math.min(i + maxLenght, theme.length()));
							if(i + maxLenght < theme.length() && theme.length() - i + maxLenght <= 16 ){
								part += "...";
							}
							message.append(part);
						}
					}
					else{
						message.append(theme);
					}
				}
			}
			pstmt.close(); result.close();
		}
		catch(SQLException e){
			e.printStackTrace();
		}
		return message;
	}
}
