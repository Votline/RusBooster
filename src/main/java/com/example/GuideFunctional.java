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
import java.util.ArrayList;
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
		guideMessage.setChatId(String.valueOf(userId));

		try(Connection conn = DriverManager.getConnection(url)){

			sql = "SELECT guide FROM guides WHERE task_id = ?";
			PreparedStatement sendPstmt = conn.prepareStatement(sql);
			sendPstmt.setInt(1, getUserTaskId(userId));

			ResultSet result = sendPstmt.executeQuery();
			

			if(result.next()){
				UserState userState = UserStateManager.getUserState(userId);
				userState.allWords = getGuide(getUserTaskId(userId));
				userState.currentPage = 0;
				guideMessage.setText(userState.allWords.get(userState.currentPage));

				InlineKeyboardMarkup guideKeyboard = new InlineKeyboardMarkup();
					InlineKeyboardButton back = new InlineKeyboardButton(); back.setCallbackData("back"); 
					InlineKeyboardButton next = new InlineKeyboardButton(); next.setCallbackData("next");
					InlineKeyboardButton allPages = new InlineKeyboardButton(); allPages.setCallbackData("toMain");
					back.setText("<"); next.setText(">"); allPages.setText("[" + (userState.currentPage+1) + "/" + (userState.allWords.size()) + "]");

				List<InlineKeyboardButton> row = Arrays.asList(back, next);
				guideKeyboard.setKeyboard(Arrays.asList(
					row, 
					Collections.singletonList(allPages)
				));
				guideMessage.setReplyMarkup(guideKeyboard);
			}			
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
	private List<String> getGuide(int taskId){
		String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/guide.db";
		String sql;

		List<String> message = new ArrayList<>();
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "SELECT guide FROM guides WHERE task_id = ?";
			
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setInt(1, taskId);

			ResultSet result = pstmt.executeQuery();
			if(result.next()){
				String guide = result.getString("guide");
				String[] themes = guide.split("!");
			
				for(String theme : themes){
					if(theme.length() > 512){
						int maxLenght = 509;
						for(int i = 0; i < theme.length(); i += maxLenght){
							if(theme.length() - 509 <= 16 || i > 0 && theme.length() - i <= 16){
								message.add(theme.substring(i));
								break;
							}
							String part = theme.substring(i, Math.min(i + maxLenght, theme.length()));
							if(i + maxLenght < theme.length()){
								part += "...";
							}
							message.add(part);
						}
					}
					else{
						message.add(theme);
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
