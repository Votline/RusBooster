package com.example;

import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;

import java.util.Calendar;
import java.util.TimeZone;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;

import java.time.ZonedDateTime;
import java.time.LocalDate;
import java.time.ZoneId;

import java.util.List;
import java.util.ArrayList;

public class NotificationSystem{
  Statistic statistic = new Statistic();
  private RusBooster bot;

  private ScheduledExecutorService scheduler;

  public NotificationSystem(RusBooster bot){
    this.bot = bot;
    scheduler = Executors.newScheduledThreadPool(1);
    scheduleNotification();
  }
  private void scheduleNotification(){
    List<Long> chatIds = statistic.getAllChatIds();
    for(long id : chatIds){
      scheduler.schedule(() -> {
       sendNotification(id);
       scheduleNotification();
      }, findPeriod(id), TimeUnit.HOURS);
    }
  }

  public void sendNotification(long id){
    try{
      SendMessage notify = new SendMessage();
      notify.setText("Не забудь выполнить задание чтобы увеличить свой рекорд!");
      notify.setChatId(String.valueOf(id));
      bot.execute(notify);
    }
    catch(TelegramApiException e){
      e.printStackTrace();
    }
  }
  public long findPeriod(long chatId){
    long period = 3;
    Calendar calendar = Calendar.getInstance();
    calendar.setTimeZone(TimeZone.getTimeZone("Europe/Moscow"));
    calendar.add(calendar.HOUR_OF_DAY, statistic.findUserOffset(chatId));
    int hour = calendar.get(Calendar.HOUR_OF_DAY);

    if(hour >= 0 && hour < 12){
      period = 6;
    }
    else if(hour >= 12 && hour < 21){
      period = 3;
    }
    else if(hour >= 21 && hour < 24){
      period = 1;
    }
    return period;
  }
}
