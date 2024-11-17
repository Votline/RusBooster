import java.sql.PreparedStatement;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.sql.Connection;
import java.sql.ResultSet;

import org.telegram.telegrambots.meta.api.objects.Update;
import org.telegram.telegrambots.meta.api.objects.Message;
import org.telegram.telegrambots.bots.TelegramLongPollingBot;
import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;

import org.telegram.telegrambots.meta.api.objects.replykeyboard.ReplyKeyboardMarkup;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.buttons.KeyboardRow;
import org.telegram.telegrambots.meta.api.objects.replykeyboard.buttons.KeyboardButton;

import java.util.List;
import java.util.ArrayList;
import java.io.IOException;

public class ReplyKeyboard{
	private MainKeyboard main = new MainKeyboard();
	private CheckKeyboard check = new CheckKeyboard();
	private SendMessage messageMenu = new SendMessage();

	private boolean isChoosing = false;
	private Settings settings = new Settings();

	public void createMenu(TelegramLongPollingBot bot, long chatId, String messageText, String userId){
		if(messageText.equals("Выбрать задание")){
			this.messageMenu = check.createMenu(chatId, userId);
			isChoosing = true;
		}
		else{
			this.messageMenu = main.createMenu(chatId);
			try{
				if(isChoosing){
 					if (Integer.parseInt(messageText) <= 26 &&  Integer.parseInt(messageText) >= 1){
						settings.chooseExercise(userId, messageText);
						System.out.println("Делегировал " + messageText + " to Settings.json");
					}
					isChoosing = false;
				}
			}
			catch(NumberFormatException e){
				e.printStackTrace();
				System.out.println("ОТКАЗ!");
			}
		}

		try{bot.execute(this.messageMenu); }
		catch(TelegramApiException e){e.printStackTrace(); }
	}

}

class MainKeyboard{
	public SendMessage createMenu(long chatId){
		SendMessage message = new SendMessage();
		message.setChatId(String.valueOf(chatId));
		message.setText("Выберите опцию: ");

		ReplyKeyboardMarkup keyboardMarkup = new ReplyKeyboardMarkup();
		keyboardMarkup.setResizeKeyboard(true);

		List<KeyboardRow> keyboard = new ArrayList<>();
		KeyboardRow row = new KeyboardRow();
		row.add(new KeyboardButton("Выбрать задание"));
		row.add(new KeyboardButton("Проверить знания"));
		row.add(new KeyboardButton("Статистика"));
		keyboard.add(row);

		keyboardMarkup.setKeyboard(keyboard);
		message.setReplyMarkup(keyboardMarkup);

		return message;
	}
}

class CheckKeyboard{
	private String url = "jdbc:sqlite:" + System.getProperty("user.dir") + "/settings.db";
	private String sql;

	public SendMessage createMenu(long chatId, String userId){
		SendMessage message = new SendMessage();
		message.setChatId(String.valueOf(chatId));
		message.setText("Напишите номер упражнения от 1 до 26: ");

		ReplyKeyboardMarkup keyboardMarkup = new ReplyKeyboardMarkup();
		keyboardMarkup.setResizeKeyboard(true);
		List<KeyboardRow> keyboard = new ArrayList<>();
		KeyboardRow row2 = new KeyboardRow();

		try(Connection conn = DriverManager.getConnection(url)){
			sql = "SELECT task_id FROM settings WHERE user_id = ?";
			PreparedStatement pstmt = conn.prepareStatement(sql);
			pstmt.setString(1, userId);
			ResultSet result = pstmt.executeQuery();
			if(result.next()){
				KeyboardRow row1 = new KeyboardRow();
				row1.add(new KeyboardButton("Текущее задание: " + result.getString("task_id") ));
				keyboard.add(row1);
			}
		}
		catch(SQLException e){
			e.printStackTrace();
		}
		row2.add(new KeyboardButton("Назад"));
		keyboard.add(row2);

		keyboardMarkup.setKeyboard(keyboard);
		message.setReplyMarkup(keyboardMarkup);

		return message;
	}
}
