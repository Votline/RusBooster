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
			this.messageMenu = check.createMenu(chatId);
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

		System.out.println(messageText);

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
	public SendMessage createMenu(long chatId){
		SendMessage message = new SendMessage();
		message.setChatId(String.valueOf(chatId));
		message.setText("Выберите задание: ");

		ReplyKeyboardMarkup keyboardMarkup = new ReplyKeyboardMarkup();
		keyboardMarkup.setResizeKeyboard(true);

		List<KeyboardRow> keyboard = new ArrayList<>();
		KeyboardRow row1 = new KeyboardRow();
		KeyboardRow row2 = new KeyboardRow();
		row1.add(new KeyboardButton("9"));
		row1.add(new KeyboardButton("10"));
		row2.add(new KeyboardButton("Назад"));
		keyboard.add(row1);
		keyboard.add(row2);

		keyboardMarkup.setKeyboard(keyboard);
		message.setReplyMarkup(keyboardMarkup);

		return message;
	}
}