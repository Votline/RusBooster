package com.example;

import java.sql.PreparedStatement;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.Statement;
import java.sql.ResultSet;

import java.time.ZonedDateTime;
import java.util.ArrayList;
import java.time.LocalDate;
import java.time.ZoneId;
import java.util.List;

import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.InlineKeyboardMarkup;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.buttons.InlineKeyboardButton;

public class Statistic{
	private String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/statistic.db";
	private String sql;

	private String message;

	public void createTable(){
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "CREATE TABLE IF NOT EXISTS statistics ("+
				"user_id LONG NOT NULL UNIQUE," +
				"chat_id LONG NOT NULL UNIQUE," +
				"current_task INTEGER DEFAULT 0," +
				"current_score INTEGER DEFAULT 0," +
				"baddest_task INTEGER DEFAULT '0'," +
				"baddest_score INTEGER DEFAULT 0," +
				"better_task INTEGER DEFAULT '0'," +
				"better_score INTEGER DEFAULT 0," +
				"last_Active_Date INTEGER  DEFAULT '0'," +
				"streak INTEGER DEFAULT '0'," +
				"timeZone INTEGER DEFAULT '0'," +
				"streakFreeze INTEGER DEFAULT '0'," +
				"PRIMARY KEY (user_id, chat_id)" +
				");";
			Statement stmt = conn.createStatement();
			stmt.execute(sql);
		}
		catch(SQLException e){
			e.printStackTrace();
		}
	}
	public String getStatistic(String userName, long userId){
		UserStateManager.getUserState(userId).isActive = false;
		createTable();
		checkStreak(userId);
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "SELECT * FROM statistics WHERE user_id = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setLong(1, userId);
			ResultSet result = pstmt.executeQuery();
			if(result.next()){
				message = "👋Привет " + userName + "!" +
					"	\nТекущее задание: " + result.getInt("current_task") +
					"	\nНаихудшая успеваимость: №" + result.getInt("baddest_task") + ", " + result.getInt("baddest_score") +
					"	\nНаилучшая успеваимость: №" + result.getInt("better_task") + ", " + result.getInt("better_score") +
					" \nТы занимаешься уже " + result.getInt("streak") + " " + getDayForm(result.getInt("streak")) + " подряд!👏" +
					" \n" + (result.getInt("timeZone") > 0 ? "Текущий часовой пояс: +" + result.getInt("timeZone") + " от МСК" : "Текущий часовой пояс: МСК"); 
			}
			else{
				message = "Ваша статистика  пуста.\nВыберите задание чтоб заполнить её!";
			}
			
			result.close(); pstmt.close();
		}
		catch(SQLException e){
			e.printStackTrace();
			message = "Таблица статистика пуста.\nВыберите задания чтоб заполнить её!";

		}
		return message;
	}
	public String chooseExercise(long userId, long chatId, String message){
		try{
			int numberOfExercise = Integer.parseInt(message);
			if(numberOfExercise <= 26 && numberOfExercise >= 1){
				try(Connection conn = DriverManager.getConnection(url)){
					sql = "INSERT INTO statistics (user_id, chat_id, current_task) VALUES (?, ?, ?) ON CONFLICT(user_id, chat_id) DO UPDATE SET current_task = excluded.current_task";
					PreparedStatement pstmt = conn.prepareStatement(sql);
					pstmt.setLong(1, userId);
					pstmt.setLong(2, chatId);
					pstmt.setInt(3, numberOfExercise);
					pstmt.executeUpdate();			
			
					pstmt.close();
				}
				return "Текущее задание: №" + numberOfExercise;
			}
			else{
				UserStateManager.getUserState(userId).isChoosing = true;
				return "Выберите задание от 1 до 26: ";
			}
		}
		catch(NumberFormatException | SQLException e){
			e.printStackTrace();
			if(!message.isEmpty()) {
				if(e instanceof NumberFormatException){
					UserStateManager.getUserState(userId).isChoosing = true;
					return "Введите номер задания числом: ";
				}
				else{
					createTable(); chooseExercise(userId, chatId, message);
					return "Произошла ошибка при работе с базой данных. Создаю новую таблицу...";
				}
			}
			return "Перенапровляемся...";
		}
	}
	public String checkStreak(long userId){
		UserState userState = UserStateManager.getUserState(userId);

		String streakMessage = "";
		int currentDate = (int) LocalDate.now(ZoneId.of("Europe/Moscow")).toEpochDay();
		int lastActiveDate = 0;
		int difference = 0;
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
			result.close(); pstmt.close();

			difference = currentDate - lastActiveDate;
			if(difference == 1 && userState.isActive || streak == 0 && userState.isActive){
				streak += 1;
				lastActiveDate = currentDate;
				streakMessage = "Стрик обновлён! Вы занимаетесь " + streak + " " + getDayForm(streak);
			}
			else if(difference > 1){
				streak = 0;
			}	

			sql = "UPDATE statistics SET last_Active_Date = ?, streak = ? WHERE user_id = ?";
			pstmt = conn.prepareStatement(sql);
			pstmt.setInt(1, lastActiveDate);
			pstmt.setInt(2, streak);
			pstmt.setLong(3, userId);      
			pstmt.executeUpdate();
			
			result.close(); pstmt.close();
    }
		catch(SQLException e){
			e.printStackTrace();  
		}
		return streakMessage;
  }
	
	private String getDayForm(int number){
		int lastDigit = number % 10;
		int lastTwoDigits = number % 100;
		if(number == 0){
			return "дней";
		}
		else if(lastDigit == 1 && lastTwoDigits != 11){
			return "день";
		}
		else if(lastDigit >= 2 && lastDigit <= 4 || !(lastTwoDigits >= 12 && lastTwoDigits <= 14)){
			return "дня";
		}
		else{
			return "дней";
		}
	}
	public List<Long> getAllChatIds(){
		ZonedDateTime moscowDate = ZonedDateTime.now(ZoneId.of("Europe/Moscow"));
		List<Long> chatIds = new ArrayList<>(); 

		try(Connection conn = DriverManager.getConnection(url)){
			sql = "SELECT user_id, chat_id, streak, last_Active_Date, timeZone FROM statistics";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			ResultSet result = pstmt.executeQuery();
			while(result.next()){
				int currentDate = (int) moscowDate.plusHours(result.getInt("timeZone")).toLocalDate().toEpochDay();
				int streak = result.getInt("streak");
				long userId = result.getLong("user_id");
				long chatId = result.getLong("chat_id");
				int lastActiveDate = result.getInt("last_Active_Date");
				if(currentDate - lastActiveDate == 1 && streak != 0){
					chatIds.add(chatId);
				}
				else if(currentDate - lastActiveDate > 1 && lastActiveDate != 0 && streak != 0){
					sql = "UPDATE statistics SET streak = ?, last_Active_Date = ? WHERE user_id = ?";
					pstmt = conn.prepareStatement(sql);
					pstmt.setInt(1, 0);
					pstmt.setInt(2, 0);
					pstmt.setLong(3, userId);
					pstmt.executeUpdate();
				}
			}
			result.close(); pstmt.close();
		}
		catch(SQLException e){
			e.printStackTrace();
		}
		return chatIds;
	}
	public int findUserOffset(long chatId){
		int timeZone = 0;
		sql = "SELECT timeZone FROM statistics WHERE chat_id = ?";
		try(Connection conn = DriverManager.getConnection(url)){
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setLong(1, chatId);
			ResultSet result = pstmt.executeQuery();
			timeZone = result.getInt("timeZone");
		}
		catch(SQLException e){
			e.printStackTrace();
		}
		return timeZone;
	}
	public SendMessage printTimeZone(long chatId){
		UserStateManager.getUserState(chatId).isSetting = true;
		SendMessage sendMessage = new SendMessage();
		sendMessage.setChatId(String.valueOf(chatId));
		sendMessage.setText(
				(findUserOffset(chatId) == 0) ?
					"Текущий часовой пояс: МСК." :
				(findUserOffset(chatId) > 0) ?
					"Текущий часовой пояс: +" + findUserOffset(chatId) + " от МСК." :
					"Текущий часовой пояс: " + findUserOffset(chatId) + " от МСК."
		); sendMessage.setText(sendMessage.getText() + "\nНапишите свою разницу во времени от МСК:");	
		InlineKeyboardMarkup setKeyboard = new InlineKeyboardMarkup();		
		InlineKeyboardButton goBack = new InlineKeyboardButton(); goBack.setText("Вернуться в главное меню");
		goBack.setCallbackData("cancelTask");
		return sendMessage;
	}
	public String setTimeZone(long userId, String messageText){
		try(Connection conn = DriverManager.getConnection(url)){
			sql = "UPDATE statistics SET timeZone = ? WHERE user_id = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setInt(1, Integer.parseInt(messageText));
			pstmt.setLong(2, userId);		
			pstmt.executeUpdate();
			messageText = ( (Integer.parseInt(messageText) > 0) ? 
					"Успешно! Ваш часовой пояс изменён на +" + messageText + " от МСК" :
					"Успешно! Ваш часовой пояс изменён на " + messageText + " от МСК");
		}
		catch(SQLException | NumberFormatException e){
			e.printStackTrace();
			messageText = "Произошла ошибка при работе с данными";
		}
		return messageText;
	}
}
